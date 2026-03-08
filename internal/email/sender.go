package email

import (
	"errors"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// Sender delivers an invoice PDF to a recipient's email address.
type Sender interface {
	SendInvoice(to string, invoice *domain.Invoice, pdfData []byte) error
}

// NoOpSender is used when SMTP is not configured.
type NoOpSender struct{}

func (n *NoOpSender) SendInvoice(_ string, _ *domain.Invoice, _ []byte) error {
	return errors.New("email not configured: set SMTP_HOST, SMTP_USER, SMTP_PASS, SMTP_FROM")
}
