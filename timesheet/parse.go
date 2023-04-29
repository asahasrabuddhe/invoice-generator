package timesheet

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"invoiceGenerator"
)

func Parse(r io.Reader, invoice *invoiceGenerator.Invoice) error {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	activeSheet := file.GetActiveSheetIndex()
	activeSheetName := file.GetSheetName(activeSheet)

	rows, err := file.GetRows(activeSheetName)
	if err != nil {
		return err
	}

	var invoiceMonth, date time.Time
	//var timesheet []float64
	var prevWeek, currentWeek, line int
	var totalHours, totalAmount float64
	for i, row := range rows {
		if i == 0 {
			invoiceMonth, err = time.ParseInLocation("Jan 2006", strings.Split(row[0], " :")[0], time.Local)
			if err != nil {
				return err
			}

			currentWeek = int(invoiceMonth.Weekday())

			invoice.Start = time.Date(invoiceMonth.Year(), invoiceMonth.Month(), 1, 0, 0, 0, 0, time.Local)
			invoice.End = invoice.Start.AddDate(0, 1, -1)

			//timesheet = make([]float64, invoice.End.Day())

			_, lastWeek := invoice.End.ISOWeek()
			_, firstWeek := invoice.Start.ISOWeek()

			// invoice will have one line per week
			weekLines := lastWeek - firstWeek

			// invoice will also include extra lines, in addition to the week lines
			invoice.Lines = make([]invoiceGenerator.Line, weekLines+len(invoice.ExtraLines))

			continue
		}
		switch len(row) {
		// in the exported sheet, if the row only has one column, it is the date
		case 1:
			date, err = time.Parse("Mon Jan 02", row[0])
			if err != nil {
				return err
			}

			date = date.AddDate(invoiceMonth.Year(), 0, 0)
			// in the exported sheet, if the row has three columns, it is the total hours logged for that day
		case 3:
			var val string
			val, err = file.CalcCellValue(activeSheetName, fmt.Sprintf("C%d", i+1))
			if err != nil {
				return err
			}

			//timesheet[date.Day()-1], err = strconv.ParseFloat(val, 64)
			var hoursWorked float64
			hoursWorked, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}

			_, week := date.ISOWeek()
			if currentWeek != week {
				prevWeek = currentWeek
				currentWeek = week
				if totalHours != 0 {
					invoice.Lines[line] = CreateLine(date, totalHours, invoice, prevWeek)
					totalAmount += invoice.Lines[line].Amount
					totalHours = 0
					line++
				}
			}

			totalHours += hoursWorked
		}
	}

	_, week := date.ISOWeek()

	prevWeek = currentWeek
	currentWeek = week
	if totalHours != 0 {
		invoice.Lines[line] = CreateLine(date, totalHours, invoice, prevWeek)
		totalAmount += invoice.Lines[line].Amount
		totalHours = 0
		line++
	}

	invoice.Total = totalAmount

	for i, extraLine := range invoice.ExtraLines {
		invoice.Lines[line+i] = invoiceGenerator.Line{
			Description: extraLine.Description,
			Amount:      extraLine.Amount,
		}

		invoice.Total += extraLine.Amount
	}

	sort.Slice(invoice.Lines, func(i, j int) bool {
		if invoice.Lines[i].StartDate.IsZero() {
			return false
		} else {
			return invoice.Lines[i].StartDate.Before(invoice.Lines[j].StartDate)
		}
	})

	return nil
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
	month := date.Format("Jan")
	year := date.Year()

	return fmt.Sprintf("%s %s %d", day, month, year)
}

func GetStartOfWeek(year, week int) time.Time {
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func GetEndOfWeek(t time.Time) time.Time {
	o := t
	for o.Weekday() != time.Sunday {
		o = o.AddDate(0, 0, 1)
		if o.Month() != t.Month() {
			return o.AddDate(0, 0, -1)
		}
	}
	return o
}

func CreateLine(thisDay time.Time, totalHours float64, invoice *invoiceGenerator.Invoice, prevWeek int) invoiceGenerator.Line {
	start := GetStartOfWeek(thisDay.Year(), prevWeek)
	end := GetEndOfWeek(start)
	daysWorked := totalHours / 8
	daysWorked = math.Round(daysWorked*100) / 100

	var description string
	if !start.Equal(end) {
		description = fmt.Sprintf(
			"%.2f days of work done in between %s and %s\n@ US$ %.2f per day",
			daysWorked, OrdinalDate(start), OrdinalDate(end), invoice.Rate,
		)
	} else {
		description = fmt.Sprintf(
			"%.2f day of work done in on %s\n@ US$ %.2f per day",
			daysWorked, OrdinalDate(start), invoice.Rate,
		)
	}

	amount := daysWorked * invoice.Rate

	return invoiceGenerator.Line{
		StartDate:   start,
		Description: description,
		Amount:      amount,
	}
}
