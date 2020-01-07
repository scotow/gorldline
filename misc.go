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

var (
	dict = map[string]string{
		" a ":   " à ",
		"Peche": "Pêche",
		"nee":   "née",
		"Burger us": "Burger US",
		//"Bar a Legumes": "Bar à Legumes",
	}
)

func init() {
	l, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	locale = l
}

func parseDate(dateString string) (time.Time, time.Time, error) {
	var startDay, endDay, year int
	var frMonth string

	n, err := fmt.Sscanf(strings.ToLower(dateString), "%s du %d au %d %s", new(string), &startDay, &endDay, &frMonth)
	if err != nil {
		return timeZero, timeZero, err
	}

	if n != 4 {
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

	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")

	v, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}

	return v
}

func isRowEmpty(row []string) bool {
	for _, cell := range row {
		if len(cell) > 0 {
			return false
		}
	}
	return true
}

func trimSheet(sheet [][]string) [][]string {
	start, end := 0, len(sheet)-1
	for i := 0; i < len(sheet); i++ {
		if !isRowEmpty(sheet[i]) {
			start = i
			break
		}
	}

	for i := len(sheet) - 1; i >= 0; i-- {
		if !isRowEmpty(sheet[i]) {
			end = i
			break
		}
	}

	return sheet[start : end+1]
}

func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return midnight(t).Add(time.Hour*24 - time.Nanosecond)
}

func smoothGrammar(s string) string {
	for k, v := range dict {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
