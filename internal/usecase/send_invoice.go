package usecase

import (
	"context"
	"fmt"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/internal/email"
	"github.com/belayhun-arage/billing-service/internal/pdf"
)

type SendInvoiceUsecase struct {
	invoiceRepo  domain.InvoiceRepository
	customerRepo domain.CustomerRepository
	emailSender  email.Sender
}

func NewSendInvoiceUsecase(
	invoiceRepo domain.InvoiceRepository,
	customerRepo domain.CustomerRepository,
	emailSender email.Sender,
) *SendInvoiceUsecase {
	return &SendInvoiceUsecase{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		emailSender:  emailSender,
	}
}

// GeneratePDF returns the raw PDF bytes for the given invoice.
func (u *SendInvoiceUsecase) GeneratePDF(ctx context.Context, invoiceID string) ([]byte, error) {
	invoice, customer, err := u.fetch(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	return pdf.GenerateInvoice(invoice, customer)
}

// Send generates a PDF invoice and emails it to the customer.
func (u *SendInvoiceUsecase) Send(ctx context.Context, invoiceID string) error {
	invoice, customer, err := u.fetch(ctx, invoiceID)
	if err != nil {
		return err
	}

	pdfData, err := pdf.GenerateInvoice(invoice, customer)
	if err != nil {
		return fmt.Errorf("generating PDF: %w", err)
	}

	return u.emailSender.SendInvoice(customer.Email, invoice, pdfData)
}

func (u *SendInvoiceUsecase) fetch(ctx context.Context, invoiceID string) (*domain.Invoice, *domain.Customer, error) {
	invoice, err := u.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, nil, fmt.Errorf("invoice %s not found: %w", invoiceID, err)
	}
	customer, err := u.customerRepo.GetByID(ctx, invoice.CustomerID)
	if err != nil {
		return nil, nil, fmt.Errorf("customer %s not found: %w", invoice.CustomerID, err)
	}
	return invoice, customer, nil
}
