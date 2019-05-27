package gorldline

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
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
	l, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	locale = l
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

type Menu struct {
	Link  string
	Start time.Time
	End   time.Time
}
