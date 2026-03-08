package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// GenerateInvoice renders an Invoice as a PDF and returns the raw bytes.
func GenerateInvoice(invoice *domain.Invoice, customer *domain.Customer) ([]byte, error) {
	f := fpdf.New("P", "mm", "A4", "")
	f.SetMargins(20, 20, 20)
	f.AddPage()

	// ── Header ─────────────────────────────────────────────────────────────
	f.SetFont("Helvetica", "B", 24)
	f.SetTextColor(30, 30, 30)
	f.Cell(0, 12, "INVOICE")
	f.Ln(14)

	// ── Invoice meta ────────────────────────────────────────────────────────
	f.SetFont("Helvetica", "", 10)
	f.SetTextColor(100, 100, 100)
	f.Cell(0, 6, fmt.Sprintf("Invoice ID : %s", invoice.ID))
	f.Ln(6)
	f.Cell(0, 6, fmt.Sprintf("Date       : %s", invoice.CreatedAt.Format("January 2, 2006")))
	f.Ln(6)
	f.Cell(0, 6, fmt.Sprintf("Status     : %s", strings.ToUpper(invoice.Status)))
	f.Ln(10)

	// ── Divider ─────────────────────────────────────────────────────────────
	f.SetDrawColor(200, 200, 200)
	f.Line(20, f.GetY(), 190, f.GetY())
	f.Ln(10)

	// ── Bill To ─────────────────────────────────────────────────────────────
	f.SetFont("Helvetica", "B", 11)
	f.SetTextColor(30, 30, 30)
	f.Cell(0, 7, "Bill To")
	f.Ln(8)
	f.SetFont("Helvetica", "", 10)
	f.SetTextColor(60, 60, 60)
	f.Cell(0, 6, customer.Name)
	f.Ln(6)
	f.Cell(0, 6, customer.Email)
	f.Ln(14)

	// ── Line items table ────────────────────────────────────────────────────
	currency := strings.ToUpper(invoice.Currency)
	amount := fmt.Sprintf("%.2f", float64(invoice.Amount)/100)

	f.SetFillColor(245, 245, 245)
	f.SetFont("Helvetica", "B", 10)
	f.SetTextColor(60, 60, 60)
	f.CellFormat(110, 8, "Description", "1", 0, "L", true, 0, "")
	f.CellFormat(30, 8, "Currency", "1", 0, "C", true, 0, "")
	f.CellFormat(30, 8, "Amount", "1", 1, "R", true, 0, "")

	f.SetFont("Helvetica", "", 10)
	f.SetFillColor(255, 255, 255)
	f.CellFormat(110, 8, "Billing Services", "1", 0, "L", false, 0, "")
	f.CellFormat(30, 8, currency, "1", 0, "C", false, 0, "")
	f.CellFormat(30, 8, amount, "1", 1, "R", false, 0, "")

	f.Ln(8)

	// ── Total ───────────────────────────────────────────────────────────────
	f.SetFont("Helvetica", "B", 11)
	f.SetTextColor(30, 30, 30)
	f.SetX(140)
	f.CellFormat(20, 8, "Total:", "", 0, "R", false, 0, "")
	f.CellFormat(30, 8, fmt.Sprintf("%s %s", currency, amount), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := f.Output(&buf); err != nil {
		return nil, fmt.Errorf("rendering PDF: %w", err)
	}
	return buf.Bytes(), nil
}
