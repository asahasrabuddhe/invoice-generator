package summary

import (
	"math"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/pdf"
)

func Action(c *cli.Context) error {
	invoices := c.StringSlice("invoice")

	sheet := Sheet{
		Month: time.Now().In(time.Local).AddDate(0, -1, 0).Format("January 2006"),
	}

	sheet.Lines = make([]Line, len(invoices))

	for i, invoice := range invoices {
		f, err := os.Open(invoice)
		if err != nil {
			return err
		}

		inv, err := pdf.Read(f)
		if err != nil {
			return err
		}

		sheet.Lines[i] = Line{
			Resource:      inv.Resource(),
			InvoiceNumber: inv.Number(),
			Amount:        inv.Amount(),
		}
	}

	// include bank fees
	sheet.Lines = append(sheet.Lines, Line{
		Resource:      "Bank fees",
		InvoiceNumber: "",
		Amount:        15.0,
	})

	for _, line := range sheet.Lines {
		sheet.Total += line.Amount
	}

	sheet.Total = math.Round(sheet.Total*100) / 100

	outFileName := "summary_statement_" + strings.ReplaceAll(strings.ToLower(sheet.Month), " ", "_") + ".pdf"

	return pdf.Generate(c.Context, "summary-statement.html.tpl", outFileName, sheet)
}
