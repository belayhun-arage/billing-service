package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Config holds all environment-driven configuration for the billing service.
// Call Load() once at startup; pass the struct to components that need it.
type Config struct {
	DatabaseURL    string
	StripeKey      string
	AllowedOrigins []string
	HTTPPort       string
	GRPCPort       string

	// SMTP — all optional; if SMTPHost is empty the email feature is disabled.
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	// Rate limiting — requests per second per API key and burst size.
	RateLimitRPS   float64
	RateLimitBurst int

	// Kafka — optional; if empty the outbox worker uses a no-op publisher.
	KafkaBrokers string
	KafkaTopic   string
}

// Load reads environment variables, applies defaults for optional fields,
// and returns an error if any required field is missing.
func Load() (*Config, error) {
	smtpPort := 587
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			smtpPort = parsed
		}
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		StripeKey:   os.Getenv("STRIPE_SECRET_KEY"),
		HTTPPort:    getEnvOrDefault("HTTP_PORT", "8080"),
		GRPCPort:    getEnvOrDefault("GRPC_PORT", "9090"),

		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: smtpPort,
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
		SMTPFrom: os.Getenv("SMTP_FROM"),

		RateLimitRPS:   getEnvFloat("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20),

		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:   getEnvOrDefault("KAFKA_TOPIC", "billing-events"),
	}

	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	if rawOrigins != "" {
		cfg.AllowedOrigins = strings.Split(rawOrigins, ",")
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	var missing []string

	if c.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if c.StripeKey == "" {
		missing = append(missing, "STRIPE_SECRET_KEY")
	}

	if len(missing) > 0 {
		return errors.New("missing required environment variables: " + strings.Join(missing, ", "))
	}
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
