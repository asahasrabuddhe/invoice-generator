package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

//go:embed invoice
var fs embed.FS

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "Ajitem Sahasrabuddhe",
			Email: "ajitem@engineering.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config-file",
			Usage: "path to the configuration file",
		},
	}

	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {
	configFilePath := c.String("config-file")

	outFilePath := filepath.Dir(configFilePath)

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	invoice, err := NewInvoice(configFile)
	if err != nil {
		return err
	}

	file, err := excelize.OpenFile("/Users/ajitem/Downloads/Apr_2023_Completed_Work_1682709355.xlsx")
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

			invoice.Lines = make([]Line, lastWeek-firstWeek+len(invoice.ExtraLines))

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
		invoice.Lines[line+i+1] = Line{
			Description: extraLine.Description,
			Amount:      extraLine.Amount,
		}

		invoice.Total += extraLine.Amount
	}

	var buf bytes.Buffer

	tpl := template.Must(
		template.
			New("invoice.html.tpl").
			Funcs(template.FuncMap{
				"formatDescription": FormatDescription,
				"formatAmount":      FormatAmount,
			}).
			ParseFS(fs, "invoice/invoice.html.tpl"),
	)

	err = tpl.Execute(&buf, invoice)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(outFilePath, GetFileName(invoice)+".html"))
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	path := LocateChrome()

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

func FormatDescription(line string) template.HTML {
	pattern := regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	if pattern.MatchString(line) {
		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			number := match[1]
			suffix := match[2]
			line = strings.ReplaceAll(line, match[0], fmt.Sprintf(`%s<span class="ordinal">%s</span>`, number, suffix))
		}
	}

	pattern = regexp.MustCompile(`@ US\$ (\d+\.\d+) per day`)
	if pattern.MatchString(line) {
		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			line = strings.ReplaceAll(line, match[0], fmt.Sprintf(`<span class="text-[0.75rem] font-light">%s</span>`, match[0]))
		}
	}

	line = strings.ReplaceAll(line, "\n", "<br>")

	return template.HTML(`<p class="text-sm text-left font-medium text-slate-700">` + line + `</p>`)
}

func FormatAmount(amount float64) string {
	return fmt.Sprintf(`US$ %.2f`, amount)
}

func GetFileName(invoice *Invoice) string {
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

func CreateLine(thisDay time.Time, totalHours float64, invoice *Invoice) Line {
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

	return Line{
		Description: description,
		Amount:      amount,
	}
}
