package gorldline

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sort"
	"time"
)

const (
	DefaultUrl = "http://restaurant-seclin.atosworldline.com"
	MenusUri   = "/WidgetPage.aspx?widgetId=35"
)

const (
	cookieContent = "portal_url=restaurant-seclin.atosworldline.com/; language=FR'"
)

var (
	ErrInvalidAPIResponse = errors.New("invalid response from the website")
)

func CurrentList() (*List, error) {
	return NewListFromUrl(DefaultUrl)
}

func NewListFromUrl(url string) (list *List, err error) {
	req, err := http.NewRequest("GET", url+MenusUri, nil)
	if err != nil {
		return
	}

	req.Header.Set("Cookie", cookieContent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = ErrInvalidAPIResponse
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	list, err = NewList(doc, url)
	if err != nil {
		return
	}

	err = res.Body.Close()
	if err != nil {
		list = nil
	}
	return
}

func NewList(doc *goquery.Document, url string) (list *List, err error) {
	links := doc.Find("#bd .main-content .section-content .content-right .ul-container ul li a")
	menus := make([]*Menu, 0, links.Length())

	links.Each(func(_ int, s *goquery.Selection) {
		menu, em := NewMenu(s, url)
		if em != nil {
			err = em
			return
		}
		menus = append(menus, menu)
	})

	if err != nil {
		return
	}

	list = &List{Menus: menus}
	sort.Sort(list)

	if len(menus) > 0 {
		list.Start = list.Menus[0].Start
		list.End = list.Menus[len(list.Menus)-1].End
	}

	return
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
