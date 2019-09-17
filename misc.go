package gorldline

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrCannotParseDay = errors.New("error while parsing menu day string")
)

var (
	months = [...]string{
		"janvier",
		"fevrier",
		"mars",
		"avril",
		"mai",
		"juin",
		"juillet",
		"aout",
		"septembre",
		"octobre",
		"novembre",
		"decembre",
	}
	locale   *time.Location
	timeZero = time.Time{}
)

func init() {
	l, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	locale = l
}

func parseDate(s string) (time.Time, time.Time, error) {
	var startDay, endDay, year int
	var frMonth string

	n, err := fmt.Sscanf(s, "semaine du %d au %d %s", &startDay, &endDay, &frMonth)
	if err != nil {
		return timeZero, timeZero, err
	}

	if n != 3 {
		return timeZero, timeZero, ErrCannotParseDay
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
		return timeZero, timeZero, ErrCannotParseDay
	}

	now := time.Now().In(locale)

	if endMonth == 12 && now.Month() == time.January {
		year = now.Year() - 1
	} else {
		year = now.Year()
	}

	var start time.Time
	end := time.Date(year, time.Month(endMonth), endDay, 23, 59, 59, 1e9-1, locale)

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

func parsePrice(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "CJ" {
		return -1
	}

	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ",", "", -1)

	v, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}

	return v
}
