package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"

	"github.com/belayhun-arage/billing-service/configs"
	billingv1 "github.com/belayhun-arage/billing-service/gen/billing/v1"
	grpcdelivery "github.com/belayhun-arage/billing-service/internal/delivery/grpc"
	httpdelivery "github.com/belayhun-arage/billing-service/internal/delivery/http"
	"github.com/belayhun-arage/billing-service/internal/email"
	"github.com/belayhun-arage/billing-service/internal/messaging"
	"github.com/belayhun-arage/billing-service/internal/repository/postgres"
	"github.com/belayhun-arage/billing-service/internal/usecase"
	"github.com/belayhun-arage/billing-service/internal/worker"
	"github.com/belayhun-arage/billing-service/pkg/auth"
	"github.com/belayhun-arage/billing-service/pkg/db"
	"github.com/belayhun-arage/billing-service/pkg/db/middleware"
	grpcpkg "github.com/belayhun-arage/billing-service/pkg/grpc"
	stripe "github.com/belayhun-arage/billing-service/external/stripe"

	"github.com/gin-gonic/gin"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Overload(findEnvFile()); err != nil {
		logger.Info(".env file not found, relying on environment variables")
	}

	cfg, err := configs.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	pool, err := db.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database connection established")

	stripeClient := stripe.NewStripeClient(cfg.StripeKey)

	customerRepo := postgres.NewCustomerRepository(pool)
	invoiceRepo := postgres.NewInvoiceRepository(pool)
	paymentRepo := postgres.NewPaymentRepository(pool)
	ledgerRepo := postgres.NewLedgerRepository(pool)
	outboxRepo := postgres.NewOutboxRepository(pool)
	idempotencyRepo := postgres.NewIdempotencyRepository(pool)
	apiKeyRepo := postgres.NewAPIKeyRepository(pool)

	// --- Email sender ---
	var emailSender email.Sender
	if cfg.SMTPHost != "" {
		emailSender = email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
		logger.Info("SMTP configured", "host", cfg.SMTPHost, "port", cfg.SMTPPort)
	} else {
		emailSender = &email.NoOpSender{}
		logger.Info("SMTP not configured — email sending disabled")
	}

	createCustomerUC := usecase.NewCreateCustomerUsecase(customerRepo)
	createInvoiceUC := usecase.NewCreateInvoiceUsecase(invoiceRepo)
	sendInvoiceUC := usecase.NewSendInvoiceUsecase(invoiceRepo, customerRepo, emailSender)
	createAPIKeyUC := usecase.NewCreateAPIKeyUsecase(apiKeyRepo)
	revokeAPIKeyUC := usecase.NewRevokeAPIKeyUsecase(apiKeyRepo)
	processPaymentUC := usecase.NewProcessPaymentUsecase(
		pool,
		customerRepo,
		invoiceRepo,
		paymentRepo,
		ledgerRepo,
		outboxRepo,
		stripeClient,
	)

	// --- HTTP server ---
	customerHandler := httpdelivery.NewCustomerHandler(createCustomerUC, logger)
	invoiceHandler := httpdelivery.NewInvoiceHandler(createInvoiceUC, sendInvoiceUC, logger)
	paymentHandler := httpdelivery.NewPaymentHandler(processPaymentUC, logger)
	apiKeyHandler := httpdelivery.NewAPIKeyHandler(createAPIKeyUC, revokeAPIKeyUC, logger)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Idempotency-Key", "X-API-Key", "X-Timestamp", "X-Signature"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public — issue and revoke API keys
	r.POST("/api-keys", apiKeyHandler.Create)
	r.DELETE("/api-keys/:key", apiKeyHandler.Revoke)

	// Protected — all business routes require valid API key + HMAC signature
	rateLimiter := auth.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	protected := r.Group("/")
	protected.Use(auth.HMACAuth(apiKeyRepo))
	protected.Use(auth.RateLimit(rateLimiter))
	protected.Use(middleware.IdempotencyMiddleware(idempotencyRepo))
	{
		protected.POST("/customers", customerHandler.CreateCustomer)
		protected.POST("/invoices", invoiceHandler.CreateInvoice)
		protected.GET("/invoices/:id/pdf", invoiceHandler.DownloadPDF)
		protected.POST("/invoices/:id/send", invoiceHandler.SendByEmail)
		protected.POST("/payments", paymentHandler.ProcessPayment)
	}

	// --- Shutdown context (cancelled on SIGINT / SIGTERM) ---
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Outbox worker ---
	var publisher messaging.EventPublisher
	if cfg.KafkaBrokers != "" {
		publisher = messaging.NewKafkaPublisher([]string{cfg.KafkaBrokers}, cfg.KafkaTopic)
		logger.Info("Kafka publisher configured", "brokers", cfg.KafkaBrokers, "topic", cfg.KafkaTopic)
	} else {
		publisher = messaging.NewNoOpPublisher(logger)
		logger.Info("Kafka not configured — outbox events will be dropped")
	}
	outboxWorker := worker.NewOutboxWorker(pool, publisher)
	go outboxWorker.Start(ctx)

	// --- gRPC server ---
	grpcPaymentHandler := grpcdelivery.NewPaymentHandler(processPaymentUC, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcpkg.RecoveryInterceptor(logger),
			grpcpkg.LoggingInterceptor(logger),
		),
	)
	billingv1.RegisterBillingServiceServer(grpcServer, grpcPaymentHandler)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("failed to start gRPC listener", "error", err)
		os.Exit(1)
	}

	// --- HTTP server ---
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	errCh := make(chan error, 2)

	go func() {
		logger.Info("gRPC server listening", "addr", ":"+cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	go func() {
		logger.Info("HTTP server listening", "addr", ":"+cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// --- Block until signal or fatal error ---
	select {
	case err := <-errCh:
		logger.Error("server error", "error", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	// --- Graceful shutdown (10 s drain window) ---
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}
	grpcServer.GracefulStop()
	pool.Close()

	logger.Info("shutdown complete")
}

// findEnvFile walks up from the current working directory until it finds a
// .env file. This allows running the binary from any subdirectory of the project.
func findEnvFile() string {
	dir, _ := os.Getwd()
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".env"
}
