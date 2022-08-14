package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}
type Config struct {
	Addr      string
	StaticDir string
}

func main() {

	cfg := new(Config)
	//If not passed, listen on standard port 4000
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directo
	// Note that the path given to the http.Dir function is relative to the pro
	// directory root.
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	service := &http.Server{
		Addr:     *&cfg.Addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	err := service.ListenAndServe()
	if err != nil {

		errorLog.Fatal(err)

	}

}
