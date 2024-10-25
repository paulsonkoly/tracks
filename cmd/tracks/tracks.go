package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/paulsonkoly/tracks/app"
	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db := openDB()
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	app := app.NewApp(logger, db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", viewTracks)
	mux.HandleFunc("GET /track/", viewTrack)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	err = http.ListenAndServe("0.0.0.0:9999", app.StandardChain().Then(mux))
	if err != nil {
		panic(err)
	}
}

func openDB() *sql.DB {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		panic("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dburl)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
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
	if err != nil {
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
