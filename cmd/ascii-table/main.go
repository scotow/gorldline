package main

import (
	"log"
	"os"

	"github.com/scotow/gorldline"
)

func main() {
	list, err := gorldline.CurrentList()
	if err != nil {
		log.Println(err)
		return
	}

	nearestWeek := list.Current()
	if nearestWeek == nil {
		log.Println("no menu available")
		return
	}

	nearestDay, err := nearestWeek.Nearest()
	if err != nil {
		log.Println(err)
		return
	}

	nearestDay.WriteAsciiTable(os.Stdout)
}