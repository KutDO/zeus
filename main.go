package main

import (
	"html/template"
	"log"
	"net/http"

	_ "Zeus/docs" // swagger generated docs

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Download Service API
// @version 1.0
// @description This is a service for downloading files with progress tracking.
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	if err := LoadConfig("config.toml"); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Verify paths are not empty
	if config.Paths.DownloadDirectory == "" || config.Paths.WebDirectory == "" {
		log.Fatalf("Download directories must be specified in the configuration")
	}

	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/downloads", CreateDownloadHandler).Methods("POST")
	r.HandleFunc("/progress/{id}", ProgressHandler).Methods("GET")

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	go StartWorker()

	srv := &http.Server{
		Handler:      r,
		Addr:         config.GetAddress(),
		WriteTimeout: config.GetWriteTimeout(),
		ReadTimeout:  config.GetReadTimeout(),
	}

	log.Printf("Starting server on %s", config.GetAddress())
	log.Fatal(srv.ListenAndServe())
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
