package main

import (
	"database/sql"
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
	"github.com/paulsonkoly/tracks/repository/sqlc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db := openDB()
	defer db.Close()

	queries := sqlc.New(db)
	repo := repository.New(queries, db)

	app := app.New(slog.New(), &repo, scs.New(db), template.New())

	handlers := handler.New(app)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.ViewTracks)

	mux.HandleFunc("GET /track/{id}", handlers.ViewTrack)
	mux.HandleFunc("GET /track/{id}/points", handlers.ViewTrackPoints)

	// GPX file
	mux.Handle("GET /gpxfiles", app.RequiresLogIn(http.HandlerFunc(handlers.GPXFiles)))
	mux.Handle("POST /gpxfile/upload", app.RequiresLogIn(http.HandlerFunc(handlers.PostUploadGPXFile)))
	mux.Handle("POST /gpxfile/{id}/delete", app.RequiresLogIn(http.HandlerFunc(handlers.DeleteGPXFile)))

	// User
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

	err = http.ListenAndServe(serverAddr(), app.StandardChain().Then(mux))
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

func serverAddr() string {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		return "0.0.0.0:9999"
	}
	return serverAddr
}
