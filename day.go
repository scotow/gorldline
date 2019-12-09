package gorldline

import (
	"errors"
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

type Meal struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
