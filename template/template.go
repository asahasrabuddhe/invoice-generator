package template

import (
	"embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed invoice
var fs embed.FS

func Get(name string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(template.FuncMap{
			"formatDescription": FormatDescription,
			"formatAmountInd":   FormatAmount(language.Hindi),
			"formatAmountCa":    FormatAmount(language.English),
			"add":               Add,
			"mul":               Multiply,
		}).
		ParseFS(fs, "invoice/"+name)
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

func FormatAmount(tag language.Tag) func(currency string, amount float64) string {
	printer := message.NewPrinter(tag)
	return func(currency string, amount float64) string {
		amt := printer.Sprintf(`%.2f`, amount)
		for i := len(amt); i < 8; i++ {
			amt = " " + amt
		}
		return currency + ` ` + amt
	}
}

func Add(a, b float64) float64 {
	return a + b
}

func Multiply(a, b float64) float64 {
	return a * b
}
