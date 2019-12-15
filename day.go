package gorldline

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidDayData = errors.New("invalid day data")
)

func NewDayRaw(types, names, prices []string, start, end time.Time) (*Day, error) {
	if len(types) != len(names) || len(types) != len(prices) {
		return nil, ErrInvalidDayData
	}

	meals := make(map[string][]*Meal)

	for i, t := range types {
		m := new(Meal)
		m.Name = smoothGrammar(names[i])
		m.Price = parsePrice(prices[i])

		if others, pres := meals[t]; pres {
			meals[t] = append(others, m)
		} else {
			meals[t] = []*Meal{m}
		}
	}

	d := new(Day)
	d.Meals = meals
	d.Start = start
	d.End = end

	return d, nil
}

func NewDay(meals map[string][]*Meal, start, end time.Time) *Day {
	d := new(Day)
	d.Meals = meals
	d.Start = start
	d.End = end

	return d
}

type Day struct {
	Meals map[string][]*Meal `json:"meals"`
	Start time.Time          `json:"start"`
	End   time.Time          `json:"end"`
}

func (d *Day) FormatFr(today bool) string {
	var b strings.Builder

	if today {
		b.WriteString("Aujourd'hui, ")
	} else {
		b.WriteString("Ce jour, ")
	}

	if v, p := d.Meals["Plat du Jour"]; p && len(v) > 0 {
		b.WriteString(v[0].Name)
		b.WriteString(" sera le plat du jour. ")
	}

	if v, p := d.Meals["Trattoria"]; p && len(v) > 0 {
		b.WriteString("Le stand Trattoria vous propose ")
		b.WriteString(v[0].Name)
		b.WriteString(". ")
	}

	if v, p := d.Meals["Cuisine du Monde"]; p && len(v) > 0 {
		b.WriteString("Le plat de la cuisine du monde sera ")
		b.WriteString(v[0].Name)
		b.WriteString(". ")
	}

	if v, p := d.Meals["Bar a Legumes"]; p && len(v) > 0 {
		acc := make([]string, 0, len(d.Meals["Bar a Legumes"]))
		for _, a := range v {
			acc = append(acc, a.Name)
		}

		b.WriteString("Les accompagnements seront ")
		b.WriteString(strings.Join(acc[:cap(acc)-2], ", "))
		b.WriteString(" et ")
		b.WriteString(acc[len(acc)-1])
		b.WriteString(". ")
	}

	return b.String()
}
