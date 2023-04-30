package template

import (
	"embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

//go:embed invoice
var fs embed.FS

func Get() (*template.Template, error) {
	return template.
		New("invoice.html.tpl").
		Funcs(template.FuncMap{
			"formatDescription": FormatDescription,
			"formatAmount":      FormatAmount,
			"calculateTax":      CalculateTax,
			"add":               Add,
		}).
		ParseFS(fs, "invoice/invoice.html.tpl")
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

func CalculateTax(amount, rate float64) float64 {
	return amount * rate / 100
}

func Add(a, b float64) float64 {
	return a + b
}
