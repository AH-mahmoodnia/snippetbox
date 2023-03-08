package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	cfg      config
}

func main() {
	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile),
	}

	flag.StringVar(&app.cfg.addr, "addr", ":4000", "HTTP network adderss")
	flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&app.cfg.dsn, "dsn", "postgres://web:123@localhost/snippetbox?sslmode=disable", "Postgresql data source name")
	flag.Parse()

	db, err := app.openDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:     app.cfg.addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLog,
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func (app *application) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", app.cfg.dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
