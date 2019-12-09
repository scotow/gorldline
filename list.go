package gorldline

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultBaseUrl = "http://restaurant-seclin.atosworldline.com"
	MenusUri       = "/WidgetPage.aspx?widgetId=35"
)

const (
	cookieContent = "portal_url=restaurant-seclin.atosworldline.com/; language=FR'"
)

var (
	ErrInvalidAPIResponse = errors.New("invalid response from the website")
)

func CurrentList() (*List, error) {
	return NewListFromUrl(DefaultBaseUrl)
}

func NewListFromUrl(baseUrl string) (*List, error) {
	req, err := http.NewRequest("GET", baseUrl+MenusUri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cookieContent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidAPIResponse
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	list, err := NewList(doc, baseUrl)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func NewList(doc *goquery.Document, baseUrl string) (*List, error) {
	links := doc.Find("#bd .main-content .section-content .content-right .ul-container ul li a")
	weeks := make([]*Week, 0, links.Length())

	var err error
	links.Each(func(_ int, s *goquery.Selection) {
		week, ew := NewWeekNode(s, baseUrl)
		if ew != nil {
			err = ew
			return
		}
		weeks = append(weeks, week)
	})

	if err != nil {
		return nil, err
	}

	l := new(List)
	l.Weeks = weeks
	sort.Sort(l)

	if len(weeks) > 0 {
		l.Start = l.Weeks[0].Start
		l.End = l.Weeks[len(l.Weeks)-1].End
	}

	return l, err
}

type List struct {
	Weeks []*Week   `json:"weeks"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (l *List) Len() int {
	return len(l.Weeks)
}

func (l *List) Less(i, j int) bool {
	return l.Weeks[i].Start.Before(l.Weeks[j].Start)
}

func (l *List) Swap(i, j int) {
	l.Weeks[i], l.Weeks[j] = l.Weeks[j], l.Weeks[i]
}

func (l *List) Current() *Week {
	now := time.Now()
	for _, week := range l.Weeks {
		if now.After(week.Start) && now.Before(week.End) {
			return week
		}
	}

	return nil
}

func (l *List) Nearest() *Week {
	now := time.Now()
	for _, week := range l.Weeks {
		if (now.After(week.Start) && now.Before(week.End)) || now.Before(week.Start) {
			return week
		}
	}

	if len(l.Weeks) > 0 {
		return l.Weeks[len(l.Weeks)-1]
	}

	return nil
}

func (l *List) Merge(other *List) {
	newWeek := make([]*Week, 0)

	for _, w1 := range other.Weeks {
		for _, w2 := range l.Weeks {
			if w1.Start == w2.Start && w2.End == w2.End {
				newWeek = append(newWeek, w1)
				break
			}
		}
	}

	l.Weeks = append(l.Weeks, newWeek...)
	sort.Sort(l)
}

func (l *List) MergeWithCurrent() error {
	l2, err := CurrentList()
	if err != nil {
		return err
	}

	l.Merge(l2)
	return nil
}
