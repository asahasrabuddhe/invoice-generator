package invoice

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/pdf"
)

// regex to match the line in this format hours (in float):mmyyyy
var lineRegex = regexp.MustCompile(`^(\d*\.?\d*):(\d+)$`)

func Action(c *cli.Context) error {
	configFile, err := os.Open(c.String("config-file"))
	if err != nil {
		return err
	}

	invoice, err := NewInvoice(configFile)
	if err != nil {
		return err
	}

	if invoice.Currency == "" {
		invoice.Currency = "US$"
	}

	if invoice.Mode == "" {
		invoice.Mode = "daily"
	}

	if invoice.Layout == "" {
		invoice.Layout = "monthly"
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	outFilePath := c.String("output-file")

	var timesheetFile *os.File

	timesheetFile, err = os.Open(c.String("timesheet-path"))
	if err == nil {
		err = Parse(timesheetFile, invoice)
		if err != nil {
			return err
		}

		err = timesheetFile.Close()
		if err != nil {
			return err
		}
	}

	if err != nil && len(invoice.ExtraLines) == 0 {
		return err
	}

	lines := len(invoice.Lines)
	if lines == 0 {
		invoice.Lines = make([]Line, len(invoice.ExtraLines))
	}
	for i, extraLine := range invoice.ExtraLines {
		if splitLine := lineRegex.FindStringSubmatch(extraLine.Description); len(splitLine) > 0 {
			if len(splitLine) != 3 {
				return errors.New("invalid line format")
			}

			var month Month
			var hours float64

			hours, err = strconv.ParseFloat(splitLine[1], 64)
			if err != nil {
				return err
			}

			month.t, err = time.ParseInLocation("012006", splitLine[2], time.Local)
			if err != nil {
				return err
			}

			invoice.Lines[i] = CreateLine(month, hours, invoice)
		} else {
			invoice.Lines[lines+i] = Line{
				Description: extraLine.Description,
				Amount:      extraLine.Amount,
			}
		}
		invoice.Total += invoice.Lines[i].Amount
	}

	if outFilePath == "" {
		outFilePath = invoice.FileName()
	}

	return pdf.Generate(c.Context, "invoice.html.tpl", outFilePath, invoice)
}
