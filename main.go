package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

//go:embed invoice
var fs embed.FS

const WorkingDaysPerWeek = 5

type Invoice struct {
	InvoiceNumber int
	InvoiceDate   string
	Lines         []Line
	Total         int
}

type Line struct {
	Description string
	Amount      string
}

type Config struct {
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	InvoiceNumber int       `json:"invoiceNumber"`
	Rate          int       `json:"rate"`
	ExtraLines    []Line    `json:"extraLines"`
}

func main() {
	config := &Config{}

	configFile, _ := os.Open("config.json")

	err := json.NewDecoder(configFile).Decode(config)
	if err != nil {
		log.Fatalln(err)
	}

	invoice := GetInvoice(config)

	file, _ := os.Create("invoice/output.html")

	tpl := template.Must(template.ParseFS(fs, "invoice/invoice.html.tpl"))
	_ = tpl.Execute(file, invoice)

	_ = file.Close()
}

func GetInvoice(config *Config) Invoice {
	var workingDaysCount, workingWeeksCount int
	for date := config.StartDate; date.Unix() <= config.EndDate.Unix(); date = date.Add(24 * time.Hour) {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			workingDaysCount++
		}
		if date.Weekday() == time.Friday {
			workingWeeksCount++
		}
	}

	if config.EndDate.Weekday() != time.Friday {
		workingWeeksCount++
	}

	var lengthOfFirstWeek, lengthOfLastWeek int

	for date := config.StartDate; date.Weekday() != time.Saturday; date = date.Add(24 * time.Hour) {
		lengthOfFirstWeek++
	}

	for date := config.EndDate; date.Weekday() != time.Sunday; date = date.Add(-24 * time.Hour) {
		lengthOfLastWeek++
	}

	workingDays := make([][]time.Time, workingWeeksCount)

	for i := 0; i < workingWeeksCount; i++ {
		if i == 0 {
			workingDays[i] = make([]time.Time, lengthOfFirstWeek)
		} else if i == workingWeeksCount-1 {
			workingDays[i] = make([]time.Time, lengthOfLastWeek)
		} else {
			workingDays[i] = make([]time.Time, WorkingDaysPerWeek)
		}
	}

	var count int
	for i := 0; i < len(workingDays); i++ {
		for j := 0; j < len(workingDays[i]); j++ {
			workingDays[i][j] = config.StartDate.Add(time.Duration(count) * 24 * time.Hour)
			count++
		}
	}

	var invoice Invoice

	invoice.Lines = make([]Line, 10)

	var total int
	for i := 0; i < len(workingDays); i++ {
		start := workingDays[i][0]
		end := workingDays[i][len(workingDays[i])-1]

		unit := "days"
		days := int(end.Sub(start).Hours()/24) + 1
		if days == 1 {
			unit = "day"
		}

		value := config.Rate * days
		total += value

		invoice.Lines[i].Description = fmt.Sprintf("%d %s of work done in between %s and %s @ US%d per day", days, unit, OrdinalDate(start), OrdinalDate(end), config.Rate)
		invoice.Lines[i].Amount = fmt.Sprintf("USD %d", value)
	}

	invoice.Total = total
	invoice.InvoiceDate = config.StartDate.Format("02-01-2006")
	invoice.InvoiceNumber = config.InvoiceNumber

	return invoice
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
	month := date.Month().String()
	year := date.Year()

	return fmt.Sprintf("%s %s %d", day, month, year)
}
