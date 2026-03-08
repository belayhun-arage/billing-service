package email

import (
	"bytes"
	"fmt"
	"strings"

	gomail "github.com/wneessen/go-mail"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// SMTPSender sends emails via an SMTP server.
type SMTPSender struct {
	host string
	port int
	user string
	pass string
	from string
}

func NewSMTPSender(host string, port int, user, pass, from string) *SMTPSender {
	return &SMTPSender{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *SMTPSender) SendInvoice(to string, invoice *domain.Invoice, pdfData []byte) error {
	msg := gomail.NewMsg()

	if err := msg.From(s.from); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("invalid to address: %w", err)
	}

	msg.Subject(fmt.Sprintf("Your Invoice #%.8s", invoice.ID))
	msg.SetBodyString(gomail.TypeTextHTML, buildBody(invoice))

	if err := msg.AttachReader(
		fmt.Sprintf("invoice-%.8s.pdf", invoice.ID),
		bytes.NewReader(pdfData),
	); err != nil {
		return fmt.Errorf("attaching PDF: %w", err)
	}

	client, err := gomail.NewClient(s.host,
		gomail.WithPort(s.port),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(s.user),
		gomail.WithPassword(s.pass),
	)
	if err != nil {
		return fmt.Errorf("creating SMTP client: %w", err)
	}

	return client.DialAndSend(msg)
}

func buildBody(inv *domain.Invoice) string {
	currency := strings.ToUpper(inv.Currency)
	amount := fmt.Sprintf("%.2f", float64(inv.Amount)/100)
	return fmt.Sprintf(`<html><body style="font-family:sans-serif;color:#333">
<h2>Invoice #%.8s</h2>
<p>Thank you for your business. Please find your invoice attached.</p>
<table cellpadding="6" style="border-collapse:collapse">
  <tr><td><strong>Amount</strong></td><td>%s %s</td></tr>
  <tr><td><strong>Status</strong></td><td>%s</td></tr>
  <tr><td><strong>Date</strong></td><td>%s</td></tr>
</table>
</body></html>`,
		inv.ID, currency, amount,
		inv.Status,
		inv.CreatedAt.Format("January 2, 2006"),
	)
}
