package main

import (
	"encoding/json"
	"fmt"
	"github.com/scotow/gorldline"
	"log"
)

const (
	timeFormat = "02 Jan 2006"
)

func main() {
	list, err := gorldline.CurrentList()
	if err != nil {
		log.Fatalln(err)
		return
	}

	nearest := list.Nearest()
	if nearest == nil {
		fmt.Println("No menu available")
		return
	}

	days, err := nearest.GetDays()
	fmt.Printf("Menu from %s to %s:\n", nearest.Start.Format(timeFormat), nearest.End.Format(timeFormat))

	data, _ := json.MarshalIndent(days, "", "\t")
	fmt.Println(string(data))
}
