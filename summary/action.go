package summary

import (
	"math"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/invoice"
	"invoiceGenerator/pdf"
)

func Action(c *cli.Context) error {
	invoices := c.StringSlice("invoice")

	sheet := Sheet{
		Month: time.Now().In(time.Local).AddDate(0, -1, 0).Format("January 2006"),
	}

	sheet.Lines = make([]Line, len(invoices))

	for i, inv := range invoices {
		f, err := os.Open(inv)
		if err != nil {
			return err
		}

		in, err := pdf.Read(f)
		if err != nil {
			return err
		}

		sheet.Lines[i] = Line{
			Resource:      in.Resource(),
			InvoiceNumber: in.Number(),
			Amount:        in.Amount(),
		}
	}

	for _, line := range sheet.Lines {
		sheet.Total += line.Amount
	}

	sheet.Total = math.Round(sheet.Total*100) / 100

	outFileName := "summary_statement_" + strings.ReplaceAll(strings.ToLower(sheet.Month), " ", "_") + ".pdf"

	err := pdf.Generate(c.Context, "summary-statement.html.tpl", outFileName, sheet)
	if err != nil {
		return err
	}

	// generate invoice inr
	configFile, err := os.Open(c.String("config-file"))
	if err != nil {
		return err
	}

	invoiceInr, err := invoice.NewInvoice(configFile)
	if err != nil {
		return err
	}

	if invoiceInr.Currency == "" {
		invoiceInr.Currency = "US$"
	}

	if invoiceInr.Mode == "" {
		invoiceInr.Mode = "daily"
	}

	if invoiceInr.Layout == "" {
		invoiceInr.Layout = "monthly"
	}

	invoiceInr.Start = time.Now().In(time.Local).AddDate(0, -1, 0)

	err = configFile.Close()
	if err != nil {
		return err
	}

	invoiceInr.Lines = []invoice.Line{
		{
			Description: "Summary Statement for " + sheet.Month,
			Amount:      sheet.Total,
		},
	}

	invoiceInr.CurrencyRate = 82.5
	invoiceInr.Total = sheet.Total

	outFilePath := invoiceInr.Number + " - " + sheet.Month + " - INR.pdf"

	return pdf.Generate(c.Context, "invoice.html.tpl", outFilePath, invoiceInr)
}
