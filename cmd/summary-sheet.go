package main

import (
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

	"invoiceGenerator/chrome"
	"invoiceGenerator/template"
)

func init() {
	err := license.SetMeteredKey(os.Getenv("UNIDOC_LICENSE_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
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

	tpl, err := template.Get("summary-sheet.html.tpl")
	if err != nil {
		return err
	}

	// get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	outFileName := "summary_sheet_" + strings.ReplaceAll(strings.ToLower(sheet.Month), " ", "_") + ".pdf"
	f, err := os.Create(filepath.Join(wd, outFileName+".html"))
	if err != nil {
		return err
	}

	err = tpl.Execute(f, sheet)
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
