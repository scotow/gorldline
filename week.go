package gorldline

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	label := s.Text()
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

func NewWeekFile(path string, start, end time.Time) (*Week, error) {
	w := new(Week)
	w.Start = start
	w.End = end
	w.LinkOrPath = path

	w.daysFetcher = func() ([]*Day, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = file.Close()
		}()

		return daysFromReader(file, start)
	}

	return w, nil
}

func NewWeekUrl(url string, start, end time.Time) (*Week, error) {
	w := new(Week)
	w.Start = start
	w.End = end
	w.LinkOrPath = url

	w.daysFetcher = func() ([]*Day, error) {
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

		return daysFromReader(bytes.NewReader(data), start)
	}

	return w, nil
}

type Week struct {
	daysFetcher func() ([]*Day, error)
	days        []*Day

	Start      time.Time
	End        time.Time
	LinkOrPath string
}

func (w *Week) GetDays() ([]*Day, error) {
	if w.days != nil {
		return w.days, nil
	}

	days, err := w.daysFetcher()
	if err != nil {
		return nil, err
	}

	w.days = days
	return days, nil
}

func daysFromReader(rc io.ReadSeeker, start time.Time) ([]*Day, error) {
	book, err := xls.OpenReader(rc, "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := trimSheet(book.ReadAllCells(64))

	// Check height.
	if len(sheet) < 4 {
		return nil, ErrInvalidSheetData
	}

	if len(sheet[0]) < 3 || (len(sheet[0])+1)%2 != 0 {
		return nil, ErrInvalidSheetData
	}

	days, err := parseDays(sheet, start)
	if err != nil {
		return nil, err
	}

	return days, nil
}

func parseDays(sheet [][]string, start time.Time) ([]*Day, error) {
	types := parseTypes(sheet)

	dayCount := (len(sheet[0]) - 1) / 2
	meals := make([][]*Meal, dayCount)
	for i := range meals {
		meals[i] = make([]*Meal, 0)
	}

	for y, row := range sheet[4:] {
		for x := 1; x < len(row); x += 2 {
			name := strings.TrimSpace(row[x])
			if name == "" {
				continue
			}

			m := new(Meal)
			m.Type = types[y]
			m.Name = name
			m.Price = parsePrice(row[x+1])

			meals[x/2] = append(meals[x/2], m)
		}
	}

	days := make([]*Day, dayCount)
	for i, mealList := range meals {
		dayStart := start.Add(time.Hour * 24 * time.Duration(i))
		days[i] = NewDay(mealList, dayStart, endOfDay(dayStart))
	}

	return days, nil
}

func parseTypes(sheet [][]string) []string {
	types := make([]string, len(sheet)-4)
	for i, row := range sheet[4:] {
		types[i] = row[0]
	}

	return types
}
