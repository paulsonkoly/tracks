package main

import (
	"html/template"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", viewTracks)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		panic(err)
	}
}

func viewTracks(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles("ui/html/base.html", "ui/html/track/track.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}
