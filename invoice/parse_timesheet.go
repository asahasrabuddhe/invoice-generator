package invoice

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func Parse(r io.Reader, in *Invoice) error {
	// open file
	file, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	// get the active sheet
	activeSheet := file.GetActiveSheetIndex()
	activeSheetName := file.GetSheetName(activeSheet)

	// get the rows from the active sheet
	rows, err := file.GetRows(activeSheetName)
	if err != nil {
		return err
	}

	var month *Month
	var day time.Time

	timesheet := make(map[int64]float64)

	// iterate over the rows
	for i, row := range rows {
		// the first row is the month. we will use this row to set things up
		if i == 0 {
			// the month of the in
			month, err = NewInvoiceMonth(strings.TrimSpace(strings.Split(row[0], ":")[0]))
			if err != nil {
				return err
			}

			in.Start = month.FirstDay()
			in.End = month.LastDay()

			continue
		}
		switch len(row) {
		// in the exported sheet, if the row only has one column, it is the day
		case 1:
			// read the day
			day, err = time.Parse("Mon Jan 02", row[0])
			if err != nil {
				return err
			}

			// set the year
			day = day.AddDate(month.Year(), 0, 0)

		// in the exported sheet, if the row has three columns, it is the total hours logged for that day
		case 3, 4:
			// the hours are calculated using the SUM formula. we need to get the value of the cell
			var val string
			val, err = file.CalcCellValue(activeSheetName, fmt.Sprintf("%c%d", len(row)+64, i+1))
			if err != nil {
				return err
			}

			// convert the hours to float64
			var hoursWorked float64
			hoursWorked, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}

			timesheet[day.Unix()] = hoursWorked
		}
	}

	var line int
	var totalHours float64

	if in.Layout == "weekly" {
		weeks := month.GetWeeks()
		in.Lines = make([]Line, len(weeks)+len(in.ExtraLines))
	} else {
		in.Lines = make([]Line, 1+len(in.ExtraLines))
	}

	for _, week := range month.GetWeeks() {
		days := int(week.End().Sub(week.Start()).Hours()) / 24

		for i := 0; i <= days; i++ {
			currentDay := week.Start().AddDate(0, 0, i)
			if hours, ok := timesheet[currentDay.Unix()]; ok {
				totalHours += hours
			}
		}

		if in.Layout == "weekly" {
			in.Lines[line] = CreateLine(week, totalHours, in)
			in.Total += in.Lines[line].Amount
			line++
		}
		in.TotalHours += totalHours
		totalHours = 0
	}
	if in.Layout == "monthly" {
		in.Lines[0] = CreateLine(month, in.TotalHours, in)
		in.Total = in.Lines[0].Amount
		line++
	}

	// sort lines with respect to date
	sort.Slice(in.Lines, func(i, j int) bool {
		if in.Lines[i].StartDate.IsZero() {
			return false
		} else {
			return in.Lines[i].StartDate.Before(in.Lines[j].StartDate)
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

func CreateLine(workPeriod Range, totalHours float64, in *Invoice) Line {
	if in.Mode == "daily" {
		daysWorked := totalHours / 8
		daysWorked = math.Round(daysWorked*100) / 100

		var description string
		if !workPeriod.Start().Equal(workPeriod.End()) {
			description = fmt.Sprintf(
				"%.2f days of work done in between %s and %s\n@ %s %.2f per day",
				daysWorked, OrdinalDate(workPeriod.Start()), OrdinalDate(workPeriod.End()), in.Currency, in.Rate,
			)
		} else {
			description = fmt.Sprintf(
				"%.2f day of work done in on %s\n@ %s %.2f per day",
				daysWorked, OrdinalDate(workPeriod.Start()), in.Currency, in.Rate,
			)
		}

		amount := daysWorked * in.Rate

		return Line{
			StartDate:   workPeriod.Start(),
			Description: description,
			Amount:      amount,
		}
	} else {
		var description string
		if !workPeriod.Start().Equal(workPeriod.End()) {
			description = fmt.Sprintf(
				"%.2f hours of work done in between %s and %s\n@ %s %.2f per hour",
				totalHours, OrdinalDate(workPeriod.Start()), OrdinalDate(workPeriod.End()), in.Currency, in.Rate,
			)
		} else {
			description = fmt.Sprintf(
				"%.2f hours of work done in on %s\n@ %s %.2f per hour",
				totalHours, OrdinalDate(workPeriod.Start()), in.Currency, in.Rate,
			)
		}

		amount := totalHours * in.Rate

		return Line{
			StartDate:   workPeriod.Start(),
			Description: description,
			Amount:      amount,
		}
	}
}

type Range interface {
	Start() time.Time
	End() time.Time
}

type Month struct {
	t time.Time
}

func (m Month) Start() time.Time {
	return time.Date(m.t.Year(), m.t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func (m Month) End() time.Time {
	return m.Start().AddDate(0, 1, -1)
}

func NewInvoiceMonth(month string) (*Month, error) {
	t, err := time.ParseInLocation("Jan 2006", month, time.Local)
	if err != nil {
		return nil, err
	}

	return &Month{t: t}, nil
}

func (m Month) Year() int {
	return m.t.Year()
}

func (m Month) FirstDay() time.Time {
	fd := time.Date(m.t.Year(), m.t.Month(), 1, 0, 0, 0, 0, time.Local)
	return fd
}

func (m Month) LastDay() time.Time {
	ld := m.FirstDay().AddDate(0, 1, -1)
	return ld
}

func (m Month) GetWeeks() []*Week {
	fd := m.FirstDay()
	ld := m.LastDay()

	_, fw := fd.ISOWeek()
	_, lw := ld.ISOWeek()
	weekCount := lw - fw + 1

	weeks := make([]*Week, weekCount)

	for i := 0; i < weekCount; i++ {
		w := NewWeek(m.t.Year(), fw+i)
		for w.start.Month() != m.t.Month() {
			w.start = w.start.AddDate(0, 0, 1)
		}
		weeks[i] = w
	}

	return weeks
}

type Week struct {
	Number int
	start  time.Time
	end    time.Time
	Hours  float64
}

func (w Week) Start() time.Time {
	return w.start
}

func (w Week) End() time.Time {
	return w.end
}

func NewWeek(year, week int) *Week {
	// Start from the middle of the year:
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

	return &Week{
		Number: week,
		start:  t,
		end:    t.AddDate(0, 0, 6),
	}
}
