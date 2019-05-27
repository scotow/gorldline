package gorldline

import (
	"errors"
	"time"
)

var (
	ErrInvalidDayData = errors.New("invalid day data")
)

func NewDay(types, names, prices []string, start, end time.Time) (*Day, error) {
	if len(types) != len(names) || len(types) != len(prices) {
		return nil, ErrInvalidDayData
	}

	meals := make([]*Meal, len(types))

	for i := range types {
		m := new(Meal)
		m.Type = types[i]
		m.Name = names[i]

		price, err := parsePrice(prices[i])
		if err != nil {
			return nil, err
		}
		m.Price = price

		meals[i] = m
	}

	d := new(Day)
	d.Meals = meals
	d.Start = start
	d.End = end

	return d, nil
}

type Day struct {
	Meals []*Meal
	Start time.Time
	End   time.Time
}
