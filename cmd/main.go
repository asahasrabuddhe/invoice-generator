package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"

	"invoiceGenerator"
	"invoiceGenerator/chrome"
	"invoiceGenerator/template"
)

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "Ajitem Sahasrabuddhe",
			Email: "ajitem.s@outlook.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config-file",
			Aliases: []string{"c"},
			Usage:   "path to the configuration file",
		},
		&cli.StringFlag{
			Name:    "timesheet-path",
			Aliases: []string{"t"},
			Usage:   "path to the timesheet file",
		},
		&cli.StringFlag{
			Name:    "output-file",
			Aliases: []string{"o"},
			Usage:   "path to the output file",
		},
	}

	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

//-t "/Users/ajitem/Downloads/Apr_2023_Completed_Work_1682709355.xlsx"

func Action(c *cli.Context) error {
	configFilePath := c.String("config-file")
	outFilePath := c.String("output-file")

	if outFilePath == "" {
		outFilePath = filepath.Dir(configFilePath)
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	invoice, err := invoiceGenerator.NewInvoice(configFile)
	if err != nil {
		return err
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	file, err := excelize.OpenFile(c.String("timesheet-path"))
	if err != nil {
		return err
	}

	activeSheet := file.GetActiveSheetIndex()
	activeSheetName := file.GetSheetName(activeSheet)

	rows, err := file.GetRows(activeSheetName)
	if err != nil {
		return err
	}

	var t, date time.Time
	var timesheet []float64
	for i, row := range rows {
		if i == 0 {
			t, err = time.ParseInLocation("Jan 2006", strings.Split(row[0], " :")[0], time.Local)
			if err != nil {
				return err
			}

			invoice.Start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
			invoice.End = invoice.Start.AddDate(0, 1, -1)

			timesheet = make([]float64, invoice.End.Day())

			_, lastWeek := invoice.End.ISOWeek()
			_, firstWeek := invoice.Start.ISOWeek()

			invoice.Lines = make([]invoiceGenerator.Line, lastWeek-firstWeek+len(invoice.ExtraLines))

			continue
		}
		switch len(row) {
		case 1:
			date, err = time.Parse("Mon Jan 02", row[0])
			if err != nil {
				return err
			}

			date = date.AddDate(t.Year(), 0, 0)
		case 3:
			var val string
			val, err = file.CalcCellValue(activeSheetName, fmt.Sprintf("C%d", i+1))
			if err != nil {
				return err
			}

			timesheet[date.Day()-1], err = strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
		}
	}

	_, currentWeek := t.ISOWeek()
	var totalHours, totalAmount float64
	var thisDay time.Time
	var line int

	for day, hours := range timesheet {
		thisDay = time.Date(t.Year(), t.Month(), day+1, 0, 0, 0, 0, time.Local)
		_, week := thisDay.ISOWeek()

		if currentWeek != week {
			currentWeek = week

			if totalHours != 0 {
				invoice.Lines[line] = CreateLine(thisDay, totalHours, invoice)
				totalHours = 0
				totalAmount += invoice.Lines[line].Amount
				line++
			}
		}

		totalHours += hours
	}

	invoice.Lines[line] = CreateLine(thisDay, totalHours, invoice)
	totalHours = 0
	totalAmount += invoice.Lines[line].Amount

	invoice.Total = totalAmount

	for i, extraLine := range invoice.ExtraLines {
		invoice.Lines[line+i+1] = invoiceGenerator.Line{
			Description: extraLine.Description,
			Amount:      extraLine.Amount,
		}

		invoice.Total += extraLine.Amount
	}

	tpl, err := template.Get()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(outFilePath, GetFileName(invoice)+".html"))
	if err != nil {
		return err
	}

	err = tpl.Execute(f, invoice)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	path := chrome.Locate()

	err = exec.
		CommandContext(
			c.Context, path, "--headless", "--disable-gpu", "--no-pdf-header-footer",
			"--print-to-pdf="+filepath.Join(outFilePath, GetFileName(invoice)),
			filepath.Join(outFilePath, GetFileName(invoice)+".html"),
		).Run()
	if err != nil {
		return err
	}

	err = os.Remove(f.Name())
	if err != nil {
		return err
	}

	return nil
}

func Ordinal(x int) string {
	var suffix string
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		} else {
			suffix = "th"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		} else {
			suffix = "th"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		} else {
			suffix = "th"
		}
	default:
		suffix = "th"
	}

	return fmt.Sprintf("%d%s", x, suffix)
}

func OrdinalDate(date time.Time) string {
	day := Ordinal(date.Day())
	month := date.Format("Jan")
	year := date.Year()

	return fmt.Sprintf("%s %s %d", day, month, year)
}

func GetFileName(invoice *invoiceGenerator.Invoice) string {
	extension := ".pdf"

	return fmt.Sprintf("%s - %s %d%s", invoice.Number, invoice.Start.Month().String(), invoice.Start.Year(), extension)
}

func GetStartOfWeek(t time.Time) time.Time {
	o := t
	for o.Weekday() != time.Monday {
		o = o.AddDate(0, 0, -1)
		if o.Month() != t.Month() {
			return o.AddDate(0, 0, 1)
		}
	}
	return o
}

func GetEndOfWeek(t time.Time) time.Time {
	o := t
	for o.Weekday() != time.Sunday {
		o = o.AddDate(0, 0, 1)
		if o.Month() != t.Month() {
			return o.AddDate(0, 0, -1)
		}
	}
	return o
}

func CreateLine(thisDay time.Time, totalHours float64, invoice *invoiceGenerator.Invoice) invoiceGenerator.Line {
	start := GetStartOfWeek(thisDay.AddDate(0, 0, -1))
	end := GetEndOfWeek(start)
	daysWorked := totalHours / 8
	daysWorked = math.Round(daysWorked*100) / 100

	var description string
	if !start.Equal(end) {
		description = fmt.Sprintf(
			"%.2f days of work done in between %s and %s\n@ US$ %.2f per day",
			daysWorked, OrdinalDate(start), OrdinalDate(end), invoice.Rate,
		)
	} else {
		description = fmt.Sprintf(
			"%.2f day of work done in on %s\n@ US$ %.2f per day",
			daysWorked, OrdinalDate(start), invoice.Rate,
		)
	}

	amount := daysWorked * invoice.Rate

	return invoiceGenerator.Line{
		Description: description,
		Amount:      amount,
	}
}
