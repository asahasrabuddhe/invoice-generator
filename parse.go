package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"strings"
	"time"
)

func ParseAction(_ *cli.Context) error {
	file, err := excelize.OpenFile("/Users/ajitem/Downloads/Apr_2023_Completed_Work_1682709355.xlsx")
	if err != nil {
		return err
	}

	activeSheet := file.GetActiveSheetIndex()
	activeSheetName := file.GetSheetName(activeSheet)

	rows, err := file.GetRows(activeSheetName)
	if err != nil {
		return err
	}

	var t, date time.Time
	var timesheet []float64
	for i, row := range rows {
		if i == 0 {
			t, err = time.Parse("Jan 2006", strings.Split(row[0], " :")[0])
			if err != nil {
				log.Fatal(err)
			}

			firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
			lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

			timesheet = make([]float64, lastOfMonth.Day())

			continue
		}
		switch len(row) {
		case 1:
			date, err = time.Parse("Mon Jan 02", row[0])
			if err != nil {
				log.Fatal(err)
			}

			date = date.AddDate(t.Year(), 0, 0)
		case 3:
			var val string
			val, err = file.CalcCellValue(activeSheetName, fmt.Sprintf("C%d", i+1))
			if err != nil {
				log.Fatal(err)
			}

			timesheet[date.Day()-1], err = strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, currentWeek := t.ISOWeek()
	var totalHours float64
	var thisDay time.Time

	for day, hours := range timesheet {
		thisDay = t.AddDate(0, 0, day)
		_, week := thisDay.ISOWeek()
		if currentWeek != week {
			currentWeek = week

			start := thisDay.AddDate(0, 0, -(7 - int(thisDay.Weekday())))
			for start.Month() != thisDay.Month() {
				start = start.AddDate(0, 0, 1)
			}
			end := start.AddDate(0, 0, 7-int(start.Weekday()))

			if totalHours != 0 {
				fmt.Printf("Between %s and %s - Days: %f - Amount - %f\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), totalHours/8, (totalHours/8)*325)
				totalHours = 0
			}
		}

		//fmt.Println(week, thisDay.Format("02 Jan 2006"), hours)
		totalHours += hours

	}

	end := thisDay.AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -6)
	for start.Month() != end.Month() {
		start = start.AddDate(0, 0, 1)
	}

	fmt.Printf("Between %s and %s - Days: %f - Amount - %f\n", start.Format("02 Jan 2006"), end.Format("02 Jan 2006"), totalHours/8, (totalHours/8)*325)

	return nil
}
