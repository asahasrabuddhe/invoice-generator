package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/urfave/cli/v2"
)

//go:embed invoice
var fs embed.FS

const WorkingDaysPerWeek = 5

type Invoice struct {
	InvoiceNumber int
	InvoiceDate   string
	FromEmail     string
	FromName      string
	FromPhone1    string
	FromPhone2    string
	FromAddress   Address
	ToName        string
	ToAddress     Address
	Lines         []Line
	Total         float64
}

type Line struct {
	Description string
	Amount      float64
}

type Config struct {
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	InvoiceNumber int       `json:"invoiceNumber"`
	Rate          float64   `json:"rate"`
	FromEmail     string    `json:"fromEmail"`
	FromName      string    `json:"fromName"`
	FromPhone1    string    `json:"fromPhone1"`
	FromPhone2    string    `json:"fromPhone2"`
	FromAddress   Address   `json:"fromAddress"`
	ToName        string    `json:"toName"`
	ToAddress     Address   `json:"toAddress"`
	ExtraLines    []Line    `json:"extraLines"`
}

type Address struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Line3 string `json:"line3"`
}

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "Ajitem Sahasrabuddhe",
			Email: "ajitem@engineering.com",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "parse-time-sheet",
			Usage:  "parse time sheet",
			Action: ParseAction,
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
		log.Fatalln(err)
	}
}

func Action(c *cli.Context) error {
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

	var buf bytes.Buffer

	tpl := template.Must(
		template.
			New("invoice.html.tpl").
			Funcs(template.FuncMap{"formatDescription": FormatDescription}).
			ParseFS(fs, "invoice/invoice.html.tpl"),
	)

	err = tpl.Execute(&buf, invoice)
	if err != nil {
		return err
	}

	page := wkhtmltopdf.NewPageReader(&buf)
	page.EnableLocalFileAccess.Set(true)
	page.EnableLocalFileAccess.Set(true)
	page.Zoom.Set(1.0)

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.Dpi.Set(600)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.AddPage(page)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	err = pdfg.WriteFile(filepath.Join(outFilePath, GetFileName(config)))
	if err != nil {
		return err
	}

	return nil
}

func GetInvoice(config Config) Invoice {
	//var workingDaysCount, workingWeeksCount int
	//for date := config.StartDate; date.Unix() <= config.EndDate.Unix(); date = date.Add(24 * time.Hour) {
	//	if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
	//		workingDaysCount++
	//	}
	//	if date.Weekday() == time.Friday {
	//		workingWeeksCount++
	//	}
	//}
	//
	//if config.EndDate.Weekday() != time.Friday {
	//	workingWeeksCount++
	//}
	//
	//var lengthOfFirstWeek, lengthOfLastWeek int
	//
	//for date := config.StartDate; date.Weekday() != time.Saturday; date = date.Add(24 * time.Hour) {
	//	lengthOfFirstWeek++
	//}
	//
	//for date := config.EndDate; date.Weekday() != time.Sunday; date = date.Add(-24 * time.Hour) {
	//	lengthOfLastWeek++
	//}
	//
	//workingDays := make([][]time.Time, workingWeeksCount)
	//
	//for i := 0; i < workingWeeksCount; i++ {
	//	if i == 0 {
	//		workingDays[i] = make([]time.Time, lengthOfFirstWeek)
	//	} else if i == workingWeeksCount-1 {
	//		workingDays[i] = make([]time.Time, lengthOfLastWeek)
	//	} else {
	//		workingDays[i] = make([]time.Time, WorkingDaysPerWeek)
	//	}
	//}
	//
	//var count int
	//for i := 0; i < len(workingDays); i++ {
	//	for j := 0; j < len(workingDays[i]); j++ {
	//		if config.StartDate.Add(time.Duration(count)*24*time.Hour).Weekday() == time.Saturday {
	//			count += 2
	//		}
	//		workingDays[i][j] = config.StartDate.Add(time.Duration(count) * 24 * time.Hour)
	//		count++
	//	}
	//}

	var invoice Invoice

	//invoice.Lines = make([]Line, 10)

	//var total int
	//for i := 0; i < len(workingDays); i++ {
	//	start := workingDays[i][0]
	//	end := workingDays[i][len(workingDays[i])-1]
	//
	//	unit := "days"
	//	days := int(end.Sub(start).Hours()/24) + 1
	//	if days == 1 {
	//		unit = "day"
	//	}
	//
	//	value := config.Rate * days
	//	total += value
	//
	//	if start == end {
	//		invoice.Lines[i].Description = fmt.Sprintf("%d %s of work done in on %s @ US%d per day", days, unit, OrdinalDate(start), config.Rate)
	//	} else {
	//		invoice.Lines[i].Description = fmt.Sprintf("%d %s of work done in between %s and %s @ US%d per day", days, unit, OrdinalDate(start), OrdinalDate(end), config.Rate)
	//	}
	//	invoice.Lines[i].Amount = value
	//}

	//invoice.Total = total
	invoice.InvoiceDate = config.StartDate.Format("02-01-2006")
	invoice.InvoiceNumber = config.InvoiceNumber
	invoice.FromEmail = config.FromEmail
	invoice.FromName = config.FromName
	invoice.FromPhone1 = config.FromPhone1
	invoice.FromPhone2 = config.FromPhone2
	invoice.FromAddress = config.FromAddress
	invoice.ToName = config.ToName
	invoice.ToAddress = config.ToAddress

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

func FormatDescription(line string) template.HTML {
	pattern := regexp.MustCompile(`(\d+)(st|nd|rd|th)`)

	if pattern.MatchString(line) {
		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			number := match[1]
			suffix := match[2]
			line = strings.Replace(line, match[0], fmt.Sprintf(`%s<span class="ft13">%s </span>`, number, suffix), -1)
		}
	}

	return template.HTML(line)
}

func FormatAmount(amount float64) string {
	return fmt.Sprintf(`USD %.2f`, amount)
}

func GetFileName(config Config) string {
	extension := ".pdf"

	return fmt.Sprintf("%d - %s %d%s", config.InvoiceNumber, config.StartDate.Month().String(), config.StartDate.Year(), extension)
}
