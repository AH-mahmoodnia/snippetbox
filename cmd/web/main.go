package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AH-mahmoodnia/snippetbox/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type config struct {
	addr string
	//staticDir string
	dsn string
}

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	cfg            config
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	app := &application{
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile),
		formDecoder: form.NewDecoder(),
	}

	flag.StringVar(&app.cfg.addr, "addr", ":4000", "HTTP network adderss")
	//flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	dsn := "postgres://web:123@localhost/snippetbox?sslmode=disable"
	flag.StringVar(&app.cfg.dsn, "dsn", dsn, "Postgresql data source name")
	flag.Parse()

	db, err := app.openDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()
	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app.sessionManager = sessionManager

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         app.cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     app.errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
