package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", viewTracks)
  mux.HandleFunc("GET /track/", viewTrack)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	err := http.ListenAndServe("0.0.0.0:9999", mux)
	if err != nil {
		panic(err)
	}
}

func viewTrack(w http.ResponseWriter, _ *http.Request) {
	gpxF, err := gpx.ParseFile("./tracks.gpx")
	if err != nil {
		panic(err)
	}

	points := []gpx.GPXPoint{}

	for _, track := range gpxF.Tracks {
		for _, segment := range track.Segments {
			points = append(points, segment.Points...)
		}
	}

  for _, route := range gpxF.Routes {
    points = append(points, route.Points...)
  }

  w.Header().Set("Content-Type", "application/json")

  err = json.NewEncoder(w).Encode(points)
	if err!= nil {
		panic(err)
	}
}

func viewTracks(w http.ResponseWriter, _ *http.Request) {


	t, err := template.ParseFiles("ui/html/base.html", "ui/html/track/track.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
