package main

import (
	"database/sql"
	"errors"
	"factcheck/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
)

var (
	l = log.New(os.Stderr, "factcheck: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	errNoAddrVar = errors.New("environment variable ADDR unset")
	errNoDbVar   = errors.New("environment variable DB_URL unset")
)

type apiConfig struct {
	DB *database.Queries
}

type config struct {
	Addr  string
	DbURL string
}

func loadEnvVars() (*config, error) {
	var c = new(config)

	envAddr, exists := os.LookupEnv("ADDR")
	if !exists {
		return nil, errNoAddrVar
	}
	c.Addr = envAddr

	envDbURL, exists := os.LookupEnv("DB_URL")
	if !exists {
		return nil, errNoDbVar
	}
	c.DbURL = envDbURL

	return c, nil
}

func main() {
	envCfg, err := loadEnvVars()
	if err != nil {
		l.Fatalf("Error parsing environment variables: %v", err)
	}
	l.Println("Parsed environment variables.")

	conn, err := sql.Open("postgres", envCfg.DbURL)
	if err != nil {
		l.Fatalf("Error opening database connection: %v", err)
	}
	l.Println("Established database connection.")

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}))
	router.Use(middleware.Logger)

	// Unauthenticated endpoints.
	router.Get("/", apiCfg.getIndex)
	router.Post("/user", apiCfg.createUserHandler)
	// Authenticated endpoints.
	router.Post("/note", apiCfg.authenticate(apiCfg.createNoteHandler))
	router.Post("/notes", apiCfg.authenticate(apiCfg.getRecentNNotesForUrl))
	router.Post("/vote", apiCfg.authenticate(apiCfg.voteOnNote))

	router.Get("/profile", apiCfg.authenticate(apiCfg.getProfile))
	l.Println("Created request router.")

	l.Printf("Starting Web service at %s.", envCfg.Addr)
	srv := &http.Server{
		Addr:    envCfg.Addr,
		Handler: router,
	}
	l.Fatal(srv.ListenAndServe())
}
