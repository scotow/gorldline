package gorldline

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

var (
	ErrInvalidLinkText = errors.New("invalid link text content")
	ErrNoLink          = errors.New("invalid menu link")
	ErrCannotParseDay  = errors.New("error while parsing menu day string")
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
	locale *time.Location
)

func init() {
	locale, _ = time.LoadLocation("Europe/Paris")
}

func NewMenu(s *goquery.Selection, baseUrl string) (*Menu, error) {
	label := strings.ToLower(s.Text())
	if label == "" {
		return nil, ErrInvalidLinkText
	}

	start, end, err := parseDate(label)
	if err != nil {
		return nil, err
	}

	uri, exists := s.Attr("href")
	if !exists || uri == "" {
		return nil, ErrNoLink
	}

	m := new(Menu)
	m.Link = baseUrl + uri
	m.Start = start
	m.End = end

	return m, nil
}

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

type Menu struct {
	Link  string
	Start time.Time
	End   time.Time
}
