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
)

var (
	invoiceNumberRegex = regexp.MustCompile(`INVOICE # (.*)`)
	titleCase          = cases.Title(language.English).String
)

type Invoice struct {
	text      string
	splitText []string
}

func NewInvoice(text string) *Invoice {
	splitText := make([]string, 0)
	for _, tt := range strings.Split(text, "\n") {
		if strings.TrimSpace(tt) != "" {
			splitText = append(splitText, tt)
		}
	}
	return &Invoice{
		text:      text,
		splitText: splitText,
	}
}

func (i Invoice) Resource() string {
	if len(i.splitText) >= 3 {
		return titleCase(i.splitText[2])
	}
	return ""
}

func (i Invoice) Number() string {
	matches := invoiceNumberRegex.FindStringSubmatch(i.text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func (i Invoice) Amount() float64 {
	amount, _ := strconv.ParseFloat(strings.Split(i.splitText[len(i.splitText)-1], " ")[1], 64)
	return amount
}

func Read(seeker io.ReadSeeker) (*Invoice, error) {
	state, _ := license.GetMeteredState()
	if !state.OK {
		err := license.SetMeteredKey(os.Getenv("UNIDOC_LICENSE_KEY"))
		if err != nil {
			return nil, err
		}
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
