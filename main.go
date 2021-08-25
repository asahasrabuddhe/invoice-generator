package main

import (
	"fmt"
	"time"
)

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

	var workingDaysCount int
	var workingWeeksCount int
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
}
