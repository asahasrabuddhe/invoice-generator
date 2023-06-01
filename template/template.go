package template

import (
	"embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"invoiceGenerator/currency"
)

//go:embed invoice
var fs embed.FS

func Get(name string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(template.FuncMap{
			"formatDescription": FormatDescription,
			"formatAmount":      FormatAmount,
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

func FormatAmount(amount float64) string {
	amt := fmt.Sprintf(`%.2f`, amount)
	for i := len(amt); i < 8; i++ {
		amt = " " + amt
	}
	return string(currency.Currency) + ` ` + amt
}
