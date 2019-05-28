package gorldline

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidLinkText  = errors.New("invalid link text content")
	ErrNoLink           = errors.New("invalid menu link")
	ErrFetchSheet       = errors.New("cannot download sheet from server")
	ErrInvalidSheetSize = errors.New("invalid sheet size")
	ErrSheetTooLarge    = errors.New("sheet is too large")
	ErrInvalidSheetData = errors.New("invalid sheet data")
)

func NewWeekNode(s *goquery.Selection, baseUrl string) (*Week, error) {
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

	return NewWeekUrl(baseUrl+uri, start, end)
}

func NewWeekUrl(url string, start, end time.Time) (*Week, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetchSheet
	}

	if resp.ContentLength <= 0 {
		return nil, ErrInvalidSheetSize
	}

	if resp.ContentLength > 1e6 {
		return nil, ErrSheetTooLarge
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	book, err := xls.OpenReader(reader, "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := book.ReadAllCells(64)

	// Check height.
	if len(sheet) < 7 {
		return nil, ErrInvalidSheetData
	}

	if len(sheet[0]) < 3 || (len(sheet[0])+1)%2 != 0 {
		return nil, ErrInvalidSheetData
	}

	days, err := parseDays(sheet)
	if err != nil {
		return nil, err
	}

	w := new(Week)
	w.Days = days
	w.Url = url
	w.Start = start
	w.End = end

	return w, nil
}

type Week struct {
	Url   string
	Days  []*Day
	Start time.Time
	End   time.Time
}

/*func (w *Week) GetDays() ([]*Day, error)  {
	if w.Days != nil {
		return w.Days, nil
	}
}*/

func parseDays(sheet [][]string) ([]*Day, error) {
	types := parseTypes(sheet)

	dayCount := (len(sheet[0]) - 1) / 2
	meals := make([][]*Meal, dayCount)
	for i := range meals {
		meals[i] = make([]*Meal, 0)
	}

	for i, row := range sheet[6:] {
		for j := 1; j < len(row); j += 2 {
			if strings.TrimSpace(row[j]) == "" {
				continue
			}

			m := new(Meal)
			m.Type = types[i]
			m.Name = row[j]
			m.Price = parsePrice(row[j+1])

			meals[j/2] = append(meals[j/2], m)
		}
	}

	days := make([]*Day, dayCount)
	for i, mealList := range meals {
		days[i] = NewDay(mealList, timeZero, timeZero)
	}

	return days, nil
}

func parseTypes(sheet [][]string) []string {
	types := make([]string, len(sheet)-6)
	for i, row := range sheet[6:] {
		types[i] = row[0]
	}

	return types
}
