package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handler"
	"github.com/paulsonkoly/tracks/app/log/slog"
	"github.com/paulsonkoly/tracks/app/session_manager/scs"
	"github.com/paulsonkoly/tracks/app/template"
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

	repo := repository.New(db)

	app := app.New(slog.New(), repo, scs.New(db), template.New())

	handlers := handler.New(app)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.ViewTrack)
	mux.HandleFunc("GET /track/", viewTrack)
	mux.Handle("GET /track/upload", app.RequiresLogIn(http.HandlerFunc(handlers.UploadTrack)))
	mux.Handle("POST /track/upload", app.RequiresLogIn(http.HandlerFunc(handlers.PostUploadTrack)))

	mux.Handle("GET /users", app.RequiresLogIn(http.HandlerFunc(handlers.ViewUsers)))
	mux.Handle("GET /user/new", app.RequiresLogIn(http.HandlerFunc(handlers.NewUser)))
	mux.Handle("POST /user/new", app.RequiresLogIn(http.HandlerFunc(handlers.PostNewUser)))
	mux.Handle("GET /user/{id}/edit", app.RequiresLogIn(http.HandlerFunc(handlers.EditUser)))
	mux.Handle("POST /user/{id}/edit", app.RequiresLogIn(http.HandlerFunc(handlers.PostEditUser)))
	mux.Handle("POST /user/{id}/delete", app.RequiresLogIn(http.HandlerFunc(handlers.DeleteUser)))
	mux.HandleFunc("GET /user/login", handlers.ViewUserLogin)
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
