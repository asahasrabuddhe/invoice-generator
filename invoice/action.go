package invoice

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/currency"
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

	switch invoice.Currency {
	case "INR":
		currency.Currency = currency.INR
	default:
		currency.Currency = currency.USD
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	if lines := c.StringSlice("lines"); len(lines) > 0 {
		invoice.Layout = "monthly"
		invoice.Mode = "hourly"
		invoice.Lines = make([]Line, len(lines))

		// regex to match the line in this format hours (in float):mmyyyy
		lineRegex := regexp.MustCompile(`^(\d*\.?\d*):(\d+)$`)
		for i, line := range lines {
			splitLine := lineRegex.FindStringSubmatch(line)

			if len(splitLine) != 3 {
				return errors.New("invalid line format")
			}

			var month Month
			var hours float64

			hours, err = strconv.ParseFloat(splitLine[1], 64)
			if err != nil {
				return err
			}

			month.t, err = time.ParseInLocation("022006", splitLine[2], time.Local)
			if err != nil {
				return err
			}

			invoice.Lines[i] = CreateLine(month, hours, invoice)
			invoice.Total += invoice.Lines[i].Amount
		}
	} else {
		var timesheetFile *os.File

		invoice.Layout = c.String("layout")
		invoice.Mode = "daily"

		timesheetFile, err = os.Open(c.String("timesheet-path"))
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
	}

	lines := len(invoice.Lines)
	for i, extraLine := range invoice.ExtraLines {
		invoice.Lines[lines+i] = Line{
			Description: extraLine.Description,
			Amount:      extraLine.Amount,
		}

		invoice.Total += extraLine.Amount
	}

	outFilePath := c.String("output-file")
	if outFilePath == "" {
		outFilePath = invoice.FileName()
	}

	return pdf.Generate(c.Context, "invoice.html.tpl", outFilePath, invoice)
}
