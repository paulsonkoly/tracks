package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handlers"
	"github.com/paulsonkoly/tracks/repository"
	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db := openDB()
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)

  repo := repository.New(db)

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	app := app.New(logger, repo, sessionManager)

  handlers := handlers.New(app)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.ViewTrack)
	mux.HandleFunc("GET /track/", viewTrack)

	mux.HandleFunc("GET /users", handlers.ViewUsers)
	mux.HandleFunc("GET /user/new", handlers.NewUser)
	// mux.HandleFunc("POST /user/new", handlers.PostNewUser)
	mux.HandleFunc("GET /user/login", viewUserLogin)
	mux.HandleFunc("POST /user/login", handlers.PostUserLogin)
	mux.HandleFunc("POST /user/logout", handlers.PostUserLogout)

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

func viewUserLogin(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles("ui/html/base.html", "ui/html/partials/navbar.html", "ui/html/user/login.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

