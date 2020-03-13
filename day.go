package gorldline

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
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

func (d *Day) FrenchSentence(today bool) string {
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

	return strings.TrimSpace(b.String())
}

func (d *Day) WriteAsciiTable(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetColWidth(12)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	data := [][]string{
		make([]string, len(d.Meals) * 2),
	}

	col := 0
	for k, v := range d.Meals {
		data[0][col] = strings.ToUpper(k)
		data[0][col + 1] = ""
		for i, m := range v {
			if i >= len(data) - 1 {
				data = append(data, make([]string, len(data[0])))
			}
			data[i + 1][col] = m.Name
			if m.Price != -1 {
				data[i + 1][col + 1] = fmt.Sprintf("%.2f€", float32(m.Price) / 100)
			}
		}
		col += 2
	}

	table.AppendBulk(data)
	table.Render()
	return nil
}

func (d *Day) WriteMarkdownTable(w io.Writer) error {
	data := [][]string{
		make([]string, len(d.Meals)),
		make([]string, len(d.Meals)),
	}

	for i := 0; i < len(d.Meals); i++ {
		data[1][i] = ":-:"
	}

	keys := make([]string, 0, len(d.Meals))
	for k := range d.Meals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	col := 0
	for _, k := range keys {
		data[0][col] = k
		v := d.Meals[k]
		for i, m := range v {
			if i >= len(data) - 2 {
				data = append(data, make([]string, len(data[0])))
			}
			if m.Price != -1 {
				data[i + 2][col] = fmt.Sprintf("%s *%.2f€*", m.Name, float32(m.Price) / 100)
			} else {
				data[i + 2][col] = m.Name
			}
		}
		col += 1
	}

	for i := 0; i < len(data); i++ {
		fmt.Printf("|%s|\n", strings.Join(data[i], "|"))
	}
	return nil
}