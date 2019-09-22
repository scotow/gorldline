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

	meals := make([]*Meal, len(types))

	for i := range types {
		m := new(Meal)
		m.Type = types[i]
		m.Name = names[i]
		m.Price = parsePrice(prices[i])

		meals[i] = m
	}

	d := new(Day)
	d.Meals = meals
	d.Start = start
	d.End = end

	return d, nil
}

func NewDay(meals []*Meal, start, end time.Time) *Day {
	d := new(Day)
	d.Meals = meals
	d.Start = start
	d.End = end

	return d
}

type Day struct {
	Meals []*Meal   `json:"meals"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Meal struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
