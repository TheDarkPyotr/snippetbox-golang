package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	config "snippetbox/config"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	database interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}

	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}

	templateCache map[string]*template.Template
}

type contextKey string

var contextKeyUser = contextKey("user")

func main() {

	cfg := config.ServerConfiguration()
	//If not passed, listen on standard port 4000
	flag.StringVar(&cfg.ServerParams.Addr, "addr", ":4000", "HTTP Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	//secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

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

	//secret: 32-byte long secret key for encrypting and authenticating
	//the session cookies
	session := sessions.New([]byte(cfg.CookieSetting.Secret32))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		database:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	service := &http.Server{
		Addr:         cfg.ServerParams.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  cfg.ServerParams.IdleTimeout,
		ReadTimeout:  cfg.ServerParams.ReadTimeout,
		WriteTimeout: cfg.ServerParams.WriteTimeout,
	}
	infoLog.Printf("Starting server on %s\n", cfg.ServerParams.Addr)

	err = service.ListenAndServeTLS(cfg.TLS.CertificatePath, cfg.TLS.PrivateKeyPath)
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
