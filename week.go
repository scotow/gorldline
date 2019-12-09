package gorldline

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
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

	Days       []*Day    `json:"days"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	LinkOrPath string    `json:"path"`
}

func (w *Week) FetchDays() error {
	days, err := w.daysFetcher()
	if err != nil {
		return err
	}

	w.Days = days
	return nil
}

func (w *Week) FetchDaysIfNeeded() error {
	if w.Days != nil {
		return nil
	}

	return w.FetchDays()
}

func (w *Week) GetDays() ([]*Day, error) {
	err := w.FetchDaysIfNeeded()
	if err != nil {
		return nil, err
	}

	return w.Days, nil
}

func (w *Week) Nearest() (*Day, error) {
	days, err := w.GetDays()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for _, d := range days {
		if now.After(d.Start) && now.Before(d.End) || now.Before(d.Start) {
			return d, nil
		}
	}

	if len(days) > 0 {
		return days[len(days)-1], nil
	}

	return nil, nil
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

	days, err := parseDays(sheet[4:], start)
	if err != nil {
		return nil, err
	}

	return days, nil
}

func parseDays(sheet [][]string, start time.Time) ([]*Day, error) {
	types := parseTypes(sheet)
	days := make([]*Day, 0, (len(sheet[0])-1)/2)

	for x, d := 1, 0; x < len(sheet[0]); x, d = x+2, d+1 {
		names := make([]string, 0, len(sheet[0]))
		prices := make([]string, 0, len(sheet[0]))
		for y := 0; y < len(sheet); y++ {
			names = append(names, strings.TrimSpace(sheet[y][x]))
			prices = append(prices, sheet[y][x+1])
		}

		dayStart := start.Add(time.Hour * 24 * time.Duration(d))
		day, err := NewDayRaw(types, names, prices, dayStart, endOfDay(dayStart))
		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}

func parseTypes(sheet [][]string) []string {
	types := make([]string, 0, len(sheet))
	for _, row := range sheet {
		types = append(types, smoothGrammar(row[0]))
	}

	return types
}
