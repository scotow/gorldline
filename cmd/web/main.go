package main

import (
	"github.com/scotow/gorldline"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	htmlTemplate = `
<!DOCTYPE html>
<html lang="en" dir="ltr">
    <head>
        <meta charset="utf-8">
        <title>Worldline - Seclin - Menu</title>
        <style media="screen">
            html, body, iframe {
                width: 100%;
                height: 100%;
                margin: 0;
                border: none;
                font-size: 0;
            }
        </style>
    </head>
    <body>
        <iframe src="https://sheet.zoho.com/sheet/view.do?url={{.Link}}&name=menu">Your browser doesn't support iFrames.</iframe>
    </body>
</html>
`
)

var (
	mainPage = template.Must(template.New("main").Parse(htmlTemplate))
)


func handle(w http.ResponseWriter, r *http.Request) {
	list, err := gorldline.CurrentList()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	nearest := list.Nearest()
	if nearest == nil {
		http.Redirect(w, r, gorldline.DefaultUrl + gorldline.MenusUri, http.StatusTemporaryRedirect)
		return
	}

	if r.RequestURI == "/direct" || r.RequestURI == "/download" {
		http.Redirect(w, r, nearest.Link, http.StatusTemporaryRedirect)
	} else {
		w.Header().Set("Content-Type", "text/html")
		_ = mainPage.Execute(w, nearest)
	}
}

func listeningAddress() string {
	port, set := os.LookupEnv("PORT")
	if !set {
		port = "8080"
	}

	return ":" + port
}

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(listeningAddress(), nil))
}
