package pdf

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"invoiceGenerator/internal/summary"
)

var (
	invoiceNumberRegex = regexp.MustCompile(`INVOICE # (.*)`)
	titleCase          = cases.Title(language.English).String
)

type Invoice struct {
	resource      string
	invoiceNumber string
	amount        float64
}

func NewInvoice(text string) *Invoice {
	splitText := make([]string, 0)
	for _, tt := range strings.Split(text, "\n") {
		if strings.TrimSpace(tt) != "" {
			splitText = append(splitText, tt)
		}
	}

	return &Invoice{
		resource:      getResource(splitText),
		invoiceNumber: getInvoiceNumber(text),
		amount:        getAmount(splitText),
	}
}

func getInvoiceNumber(text string) string {
	matches := invoiceNumberRegex.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func getAmount(splitText []string) float64 {
	amount, _ := strconv.ParseFloat(strings.Split(splitText[len(splitText)-1], " ")[1], 64)

	return amount
}

func getResource(splitText []string) string {
	return titleCase(splitText[2])
}

func (i *Invoice) Line() summary.Line {
	return summary.Line{
		Resource:      i.resource,
		InvoiceNumber: i.invoiceNumber,
		Amount:        i.amount,
	}
}

func Read(seeker io.ReadSeeker) (*Invoice, error) {
	err := license.SetMeteredKey(os.Getenv("UNIDOC_LICENSE_KEY"))
	if err != nil {
		return nil, err
	}

	reader, err := model.NewPdfReader(seeker)
	if err != nil {
		return nil, err
	}

	page, err := reader.GetPage(1)
	if err != nil {
		return nil, err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return nil, err
	}

	text, err := ex.ExtractText()
	if err != nil {
		return nil, err
	}

	return NewInvoice(text), nil
}
