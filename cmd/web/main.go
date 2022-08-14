package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	config "snippetbox/config"
	"snippetbox/pkg/models/mysql"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	database      *mysql.SnippetModel
	templateCache map[string]*template.Template
}
type Config struct {
	Addr      string
	StaticDir string
}

func main() {

	cfg := config.ServerConfiguration()
	//If not passed, listen on standard port 4000
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")

	//Build connection parameter readed from ENV and loaded in cfg
	dsn_conn_parameters := fmt.Sprintf("%s:%s@/%s", cfg.DbConf.Username, cfg.DbConf.Password, cfg.DbConf.DbName)
	dsn := flag.String("dsn", dsn_conn_parameters+"?parseTime=true", "MySQL Database connection parameters")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	//Create the template cache to avoid reloading template for every page
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		database:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	service := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg),
	}
	infoLog.Printf("Starting server on %s\n", cfg.Addr)

	err = service.ListenAndServe()
	if err != nil {

		errorLog.Fatal(err)

	}

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
