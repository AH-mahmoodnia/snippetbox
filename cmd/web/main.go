package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
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
	flag.Parse()

	srv := &http.Server{
		Addr:     app.cfg.addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLog,
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
