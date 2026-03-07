package middleware

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/belayhun-arage/billing-service/internal/repository/postgres"
)

type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func IdempotencyMiddleware(repo *postgres.IdempotencyRepository) gin.HandlerFunc {

	return func(c *gin.Context) {

		key := c.GetHeader("Idempotency-Key")

		if key == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		resp, status, err := repo.Get(ctx, key)

		if err == nil {

			c.Data(status, "application/json", resp)
			c.Abort()

			return
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}

		c.Writer = writer

		c.Next()

		repo.Save(
			ctx,
			key,
			writer.body.Bytes(),
			writer.Status(),
		)
	}
}
