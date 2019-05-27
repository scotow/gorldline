package gorldline

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sort"
	"time"
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
	menus := make([]*Menu, 0, links.Length())

	var err error
	links.Each(func(_ int, s *goquery.Selection) {
		menu, em := NewMenuNode(s, baseUrl)
		if em != nil {
			err = em
			return
		}
		menus = append(menus, menu)
	})

	if err != nil {
		return nil, err
	}

	l := new(List)
	l.Menus = menus
	sort.Sort(l)

	if len(menus) > 0 {
		l.Start = l.Menus[0].Start
		l.End = l.Menus[len(l.Menus)-1].End
	}

	return l, err
}

type List struct {
	Menus []*Menu
	Start time.Time
	End   time.Time
}

func (l *List) Len() int {
	return len(l.Menus)
}

func (l *List) Less(i, j int) bool {
	return l.Menus[i].Start.Before(l.Menus[j].Start)
}

func (l *List) Swap(i, j int) {
	l.Menus[i], l.Menus[j] = l.Menus[j], l.Menus[i]
}

func (l *List) Current() (menu *Menu) {
	now := time.Now()
	for _, menu = range l.Menus {
		if now.After(menu.Start) && now.Before(menu.End) {
			return
		}
	}

	menu = nil
	return
}

func (l *List) Nearest() (menu *Menu) {
	now := time.Now()
	for _, menu = range l.Menus {
		if (now.After(menu.Start) && now.Before(menu.End)) || now.Before(menu.Start) {
			return menu
		}
	}

	if len(l.Menus) > 0 {
		menu = l.Menus[len(l.Menus)-1]
	}

	return
}
