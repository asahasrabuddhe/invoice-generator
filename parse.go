package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseAction(c *cli.Context) error {
	configFilePath := c.String("config-file")

	outFilePath := filepath.Dir(configFilePath)

	config := Config{}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return err
	}

	invoice := GetInvoice(config)

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
			t, err = time.Parse("Jan 2006", strings.Split(row[0], " :")[0])
			if err != nil {
				log.Fatal(err)
			}

			config.StartDate = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
			config.EndDate = config.StartDate.AddDate(0, 1, -1)

			timesheet = make([]float64, config.EndDate.Day())

			_, lastWeek := config.EndDate.ISOWeek()
			_, firstWeek := config.StartDate.ISOWeek()

			invoice.Lines = make([]Line, lastWeek-firstWeek)

			continue
		}
		switch len(row) {
		case 1:
			date, err = time.Parse("Mon Jan 02", row[0])
			if err != nil {
				log.Fatal(err)
			}

			date = date.AddDate(t.Year(), 0, 0)
		case 3:
			var val string
			val, err = file.CalcCellValue(activeSheetName, fmt.Sprintf("C%d", i+1))
			if err != nil {
				log.Fatal(err)
			}

			timesheet[date.Day()-1], err = strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, currentWeek := t.ISOWeek()
	var totalHours, totalAmount float64
	var thisDay time.Time
	var daysWorked, line int

	for day, hours := range timesheet {
		thisDay = t.AddDate(0, 0, day)
		_, week := thisDay.ISOWeek()
		if currentWeek != week {
			currentWeek = week

			start := thisDay.AddDate(0, 0, -(7 - int(thisDay.Weekday())))
			for start.Month() != thisDay.Month() {
				start = start.AddDate(0, 0, 1)
			}
			end := start.AddDate(0, 0, 7-int(start.Weekday()))

			if totalHours != 0 {
				//fmt.Printf("Between %s and %s - Days: %f - Amount - %f\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), totalHours/8, (totalHours/8)*325)
				invoice.Lines[line] = Line{
					Description: fmt.Sprintf("%d days of work done in between %s and %s @ US$ %.2f per day", daysWorked, start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), config.Rate),
					Amount:      (totalHours / 8) * config.Rate,
				}
				totalHours = 0
				daysWorked = 0
				totalAmount += invoice.Lines[line].Amount
				line++
			}
		}

		//fmt.Println(week, thisDay.Format("02 Jan 2006"), hours)
		totalHours += hours
		if hours != 0 {
			daysWorked++
		}

	}

	end := thisDay.AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -6)
	for start.Month() != end.Month() {
		start = start.AddDate(0, 0, 1)
	}

	//fmt.Printf("Between %s and %s - Days: %f - Amount - %f\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), totalHours/8, (totalHours/8)*325)
	invoice.Lines[line] = Line{
		Description: fmt.Sprintf("%d days of work done in between %s and %s @ US$ %.2f per day", daysWorked, start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), config.Rate),
		Amount:      (totalHours / 8) * config.Rate,
	}
	totalHours = 0
	totalAmount += invoice.Lines[line].Amount

	invoice.Total = totalAmount
	invoice.InvoiceDate = time.Now().Format("02-01-2006")

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

	f, err := os.Create(filepath.Join(outFilePath, GetFileName(config)+".html"))
	if err != nil {
		return err
	}

	_, _ = f.Write(buf.Bytes())
	_ = f.Close()

	path := LocateChrome()

	err = exec.
		CommandContext(
			c.Context, path, "--headless", "--disable-gpu", "--no-pdf-header-footer",
			"--print-to-pdf="+filepath.Join(outFilePath, GetFileName(config)),
			filepath.Join(outFilePath, GetFileName(config)+".html"),
		).Run()
	if err != nil {
		return err
	}

	//_ = os.Remove(f.Name())

	return nil
}
