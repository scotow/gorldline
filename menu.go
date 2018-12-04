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

func NewMenu(s *goquery.Selection, url string) (menu *Menu, err error) {
	label := strings.ToLower(s.Text())
	if label == "" {
		err = ErrInvalidLinkText
		return
	}

	start, end, err := parseDate(label)
	if err != nil {
		return
	}

	link, exists := s.Attr("href")
	if !exists || link == "" {
		err = ErrNoLink
		return
	}

	menu = &Menu{
		Link:  url + link,
		Start: start,
		End:   end,
	}
	return
}

func parseDate(s string) (start, end time.Time, err error) {
	var startDay, endDay int
	var frMonth string

	n, err := fmt.Sscanf(s, "semaine du %d au %d %s", &startDay, &endDay, &frMonth)
	if err != nil {
		return
	}

	if n != 3 {
		err = ErrCannotParseDay
		return
	}

	var endMonth int
	for i, m := range months {
		if frMonth == m {
			endMonth = i + 1
			break
		}
	}

	if endMonth == 0 {
		err = ErrCannotParseDay
		return
	}

	now := time.Now().In(locale)
	end = time.Date(now.Year(), time.Month(endMonth), endDay, 0, 0, 0, 0, locale)

	if startDay < endDay {
		start = time.Date(now.Year(), time.Month(endMonth), startDay, 0, 0, 0, 0, locale)
	} else {
		startMonth := endMonth - 1
		if startMonth == 0 {
			startMonth = 12
		}

		start = time.Date(now.Year(), time.Month(startMonth), startDay, 0, 0, 0, 0, locale)
	}

	return
}

type Menu struct {
	Link  string
	Start time.Time
	End   time.Time
}
