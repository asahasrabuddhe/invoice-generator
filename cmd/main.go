package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"invoiceGenerator"
	"invoiceGenerator/chrome"
	"invoiceGenerator/template"
	"invoiceGenerator/timesheet"
)

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "Ajitem Sahasrabuddhe",
			Email: "ajitem.s@outlook.com",
		},
	}

	app.Usage = "igen gnereates invoices from monday.com timesheets"

	app.Commands = []*cli.Command{
		{
			Name:    "summary-statement",
			Aliases: []string{"ss"},
			Usage:   "generate summary statement",
			Action:  SummarySheetAction,
			Hidden:  true,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "invoice",
					Aliases: []string{"i"},
					Usage:   "path to the invoice file",
				},
			},
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

func Action(c *cli.Context) error {
	configFilePath := c.String("config-file")

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

	timesheetFile, err := os.Open(c.String("timesheet-path"))
	if err != nil {
		return err
	}

	err = timesheet.Parse(timesheetFile, invoice)
	if err != nil {
		return err
	}

	outFilePath := c.String("output-file")
	if outFilePath == "" {
		outFilePath = GetFileName(invoice)
	}

	return GeneratePDF(c.Context, "invoice.html.tpl", outFilePath, invoice)
}

func GetFileName(invoice *invoiceGenerator.Invoice) string {
	extension := ".pdf"

	return fmt.Sprintf("%s - %s %d%s", invoice.Number, invoice.Start.Month().String(), invoice.Start.Year(), extension)
}

type SummarySheet struct {
	Month string        `json:"month"`
	Lines []SummaryLine `json:"lines"`
	Total float64
}

type SummaryLine struct {
	Resource      string
	InvoiceNumber string
	Amount        float64
}

func SummarySheetAction(c *cli.Context) error {
	err := license.SetMeteredKey(os.Getenv("UNIDOC_LICENSE_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	invoices := c.StringSlice("invoice")

	sheet := SummarySheet{
		Month: time.Now().In(time.Local).AddDate(0, -1, 0).Format("January 2006"),
	}

	sheet.Lines = make([]SummaryLine, len(invoices))

	for i, invoice := range invoices {
		file, err := os.Open(invoice)
		if err != nil {
			return err
		}

		reader, err := model.NewPdfReader(file)
		if err != nil {
			return err
		}

		page, err := reader.GetPage(1)
		if err != nil {
			return err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return err
		}

		splitText := make([]string, 0)
		for _, tt := range strings.Split(text, "\n") {
			if strings.TrimSpace(tt) != "" {
				splitText = append(splitText, tt)
			}
		}

		sheet.Lines[i] = SummaryLine{
			Resource:      GetResource(splitText),
			InvoiceNumber: GetInvoiceNumber(text),
			Amount:        GetAmount(splitText),
		}
	}

	for _, line := range sheet.Lines {
		sheet.Total += line.Amount
	}

	sheet.Total = math.Round(sheet.Total*100) / 100

	outFileName := "summary_sheet_" + strings.ReplaceAll(strings.ToLower(sheet.Month), " ", "_") + ".pdf"

	return GeneratePDF(c.Context, "summary-sheet.html.tpl", outFileName, sheet)
}

var invoiceNumberRegex = regexp.MustCompile(`INVOICE # (.*)`)

func GetInvoiceNumber(text string) string {
	matches := invoiceNumberRegex.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func GetAmount(splitText []string) float64 {
	amount, _ := strconv.ParseFloat(strings.Split(splitText[len(splitText)-1], " ")[1], 64)

	return amount
}

var titleCase = cases.Title(language.English).String

func GetResource(splitText []string) string {
	return titleCase(splitText[2])
}

func GeneratePDF(ctx context.Context, templateName, outFileName string, templateData any) error {
	// get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(wd, outFileName+".html"))
	if err != nil {
		return err
	}

	tpl, err := template.Get(templateName)
	if err != nil {
		return err
	}

	err = tpl.Execute(f, templateData)
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
			ctx, path, "--headless", "--disable-gpu", "--no-pdf-header-footer",
			"--print-to-pdf="+outFileName,
			outFileName+".html",
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
