package main

import (
	"fmt"
	"time"
)

const WorkingDaysPerWeek = 5

func main() {
	now := time.Now()

	startDate := time.Date(now.Year(), now.Month(), 16, 0, 0, 0, 0, now.Location())
	if startDate.Weekday() == time.Saturday || startDate.Weekday() == time.Sunday {
		for true {
			startDate.Add(24 * time.Hour)
			if startDate.Weekday() != time.Saturday && startDate.Weekday() != time.Sunday {
				break
			}
		}
	}

	endDate := time.Date(now.Year(), now.Month()+1, 15, 0, 0, 0, 0, now.Location())

	var workingDaysCount, workingWeeksCount int
	for date := startDate; date.Unix() <= endDate.Unix(); date = date.Add(24 * time.Hour) {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			workingDaysCount++
		}
		if date.Weekday() == time.Friday {
			workingWeeksCount++
		}
	}

	if endDate.Weekday() != time.Friday {
		workingWeeksCount++
	}

	fmt.Println("Days", workingDaysCount)
	fmt.Println("Weeks", workingWeeksCount)

	var lengthOfFirstWeek, lengthOfLastWeek int

	for date := startDate; date.Weekday() != time.Saturday; date = date.Add(24 * time.Hour) {
		lengthOfFirstWeek++
	}

	for date := endDate; date.Weekday() != time.Sunday; date = date.Add(-24 * time.Hour) {
		lengthOfLastWeek++
	}

	fmt.Println("Length of first week", lengthOfFirstWeek)
	fmt.Println("Length of last week", lengthOfLastWeek)

	workingDays := make([][]time.Time, workingWeeksCount)

	for i := 0; i < workingWeeksCount; i++ {
		if i == 0 {
			workingDays[i] = make([]time.Time, lengthOfFirstWeek)
		} else if i == workingWeeksCount-1 {
			workingDays[i] = make([]time.Time, lengthOfLastWeek)
		} else {
			workingDays[i] = make([]time.Time, WorkingDaysPerWeek)
		}
	}

	var count int
	for i := 0; i < len(workingDays); i++ {
		for j := 0; j < len(workingDays[i]); j++ {
			workingDays[i][j] = startDate.Add(time.Duration(count) * 24 * time.Hour)
			count++
		}
	}

	var total int
	for i := 0; i < len(workingDays); i++ {
		start := workingDays[i][0]
		end := workingDays[i][len(workingDays[i])-1]

		unit := "days"
		days := int(end.Sub(start).Hours() / 24) + 1
		if days == 1 {
			unit = "day"
		}

		value := 263*days
		total += value

		fmt.Printf("%d %s of work done in between %s and %s %d @ US263 per day - %d\n", days, unit, start.String(), end.String(), start.Year(), value)
	}

	fmt.Println("Total", total)
}
