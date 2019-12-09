package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/scotow/gorldline"
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

	err = nearest.FetchDaysIfNeeded()
	if err != nil {
		log.Fatalln(err)
		return
	}

	data, err := json.MarshalIndent(nearest, "", "\t")
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Print(string(data))
}
