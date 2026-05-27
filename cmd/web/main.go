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

	"github.com/gin-contrib/sessions/cookie"
	_ "github.com/go-sql-driver/mysql"
	"guthub.com/eduartepaiva/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	store         cookie.Store
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "devuser:devpassword@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	secret, ok := os.LookupEnv("SESSION_SECRET")
	if !ok {
		panic("missing SESSION_SECRET from environment variables")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	store := cookie.NewStore([]byte(secret))

	app := application{
		infoLog,
		errorLog,
		&mysql.SnippetModel{DB: db},
		templateCache,
		store,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
