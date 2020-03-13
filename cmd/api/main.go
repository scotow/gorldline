package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/scotow/gorldline"
)

func handleCurrentWeek(w http.ResponseWriter, _ *http.Request) {
	list, err := gorldline.CurrentList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	nearestWeek := list.Nearest()
	if nearestWeek == nil {
		http.Error(w, "no menu available", http.StatusBadGateway)
		log.Println("no menu available")
		return
	}

	err = nearestWeek.FetchDaysIfNeeded()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	writeJson(nearestWeek, w)
}

func handleCurrentDay(w http.ResponseWriter, _ *http.Request) {
	list, err := gorldline.CurrentList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	nearestWeek := list.Nearest()
	if nearestWeek == nil {
		http.Error(w, "no menu available", http.StatusBadGateway)
		log.Println("no menu available")
		return
	}

	nearestDay, err := nearestWeek.Nearest()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	writeJson(nearestDay, w)
}

func handleCurrentDayFr(w http.ResponseWriter, _ *http.Request) {
	list, err := gorldline.CurrentList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	nearestWeek := list.Nearest()
	if nearestWeek == nil {
		http.Error(w, "no menu available", http.StatusBadGateway)
		log.Println("no menu available")
		return
	}

	nearestDay, err := nearestWeek.Nearest()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if nearestDay == nil {
		http.Error(w, "cannot find nearest day", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if nearestDay.Start.After(time.Now()) || nearestDay.End.Before(time.Now()) {
		_, _ = w.Write([]byte("Ce menu ne correspond pas au menu d'aujourd'hui.\n"))
		_, _ = w.Write([]byte(nearestDay.FrenchSentence(false)))
	} else {
		_, _ = w.Write([]byte(nearestDay.FrenchSentence(true)))
	}
}

func writeJson(element interface{}, w http.ResponseWriter) {
	data, err := json.Marshal(element)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/week/current", handleCurrentWeek)
	router.HandleFunc("/day/current/", handleCurrentDay)
	router.HandleFunc("/day/current/format/fr", handleCurrentDayFr)

	log.Println("Listening at", ":8080")
	log.Fatalln(http.ListenAndServe(":8080", router))
}
