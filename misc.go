package gorldline

import (
	"fmt"
	"strings"
	"time"
)

func parseDate(s string) (time.Time, time.Time, error) {
	var startDay, endDay, year int
	var frMonth string

	n, err := fmt.Sscanf(s, "semaine du %d au %d %s", &startDay, &endDay, &frMonth)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if n != 3 {
		return time.Time{}, time.Time{}, ErrCannotParseDay
	}

	frMonth = strings.Replace(frMonth, "é", "e", -1)
	frMonth = strings.Replace(frMonth, "û", "u", -1)

	var endMonth int
	for i, m := range months {
		if frMonth == m {
			endMonth = i + 1
			break
		}
	}

	if endMonth == 0 {
		return time.Time{}, time.Time{}, ErrCannotParseDay
	}

	now := time.Now().In(locale)

	if endMonth == 12 && now.Month() == time.January {
		year = now.Year() - 1
	} else {
		year = now.Year()
	}

	var start time.Time
	end := time.Date(year, time.Month(endMonth), endDay, 0, 0, 0, 0, locale)

	if startDay < endDay {
		start = time.Date(year, time.Month(endMonth), startDay, 0, 0, 0, 0, locale)
	} else {
		startMonth := endMonth - 1
		if startMonth == 0 {
			startMonth = 12
		}

		start = time.Date(year, time.Month(startMonth), startDay, 0, 0, 0, 0, locale)
	}

	return start, end, nil
}
