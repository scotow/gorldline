package main

import (
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

	fmt.Printf("Menu from %s to %s:\n", nearest.Start.Format(timeFormat), nearest.End.Format(timeFormat))
	fmt.Println(nearest.Link)
}
