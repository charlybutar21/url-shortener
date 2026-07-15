package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	appHttp "url-shortener/internal/delivery/http"
	"url-shortener/internal/repository/postgres"
	"url-shortener/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Setup DB Connection
	dsn := "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Failed to connect to DB: %v", err)
		log.Printf("Make sure you have started PostgreSQL via docker-compose up -d")
	} else {
		log.Printf("Successfully connected to Database!")
	}

	// 2. Setup Dependencies
	repo := postgres.NewPostgresURLRepository(db)
	uc := usecase.NewURLUsecase(repo)

	// 3. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 4. Register Routes
	appHttp.NewURLHandler(r, uc)

	// 5. Start Server
	fmt.Println("Server starting on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
