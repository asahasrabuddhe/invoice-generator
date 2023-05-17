package invoice

import (
	"os"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/pdf"
)

func Action(c *cli.Context) error {
	configFile, err := os.Open(c.String("config-file"))
	if err != nil {
		return err
	}

	invoice, err := NewInvoice(configFile)
	if err != nil {
		return err
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	invoice.Layout = c.String("layout")

	timesheetFile, err := os.Open(c.String("timesheet-path"))
	if err != nil {
		return err
	}

	err = Parse(timesheetFile, invoice)
	if err != nil {
		return err
	}

	err = timesheetFile.Close()
	if err != nil {
		return err
	}

	outFilePath := c.String("output-file")
	if outFilePath == "" {
		outFilePath = invoice.FileName()
	}

	return pdf.Generate(c.Context, "invoice.html.tpl", outFilePath, invoice)
}
