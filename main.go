package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	plataform      string
	secrete        string
	polka          string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("erro ao carregar as variaveis de ambiente: %v", err))
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("erro ao abrir conexao com o banco de dados: %v", err)
	}

	dbQueries := database.New(db)
	mux := http.NewServeMux()

	api := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		plataform:      os.Getenv("PLATAFORM"),
		secrete:        os.Getenv("SECRETE"),
		polka:          os.Getenv("POLKA"),
	}
	// GET
	mux.Handle("/app/", http.StripPrefix("/app", api.middlewareMetricsInt(http.HandlerFunc(handleRoot))))
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /admin/metrics", api.metrics)
	mux.HandleFunc("GET /api/chirps", api.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", api.getChirpsByID)

	// POST
	mux.HandleFunc("POST /admin/reset", api.reset)
	mux.HandleFunc("POST /api/chirps", api.addChirps)
	mux.HandleFunc("POST /api/users", api.addUser)
	mux.HandleFunc("POST /api/login", api.loginUser)
	mux.HandleFunc("POST /api/refresh", api.validateRefreshToken)
	mux.HandleFunc("POST /api/revoke", api.revokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", api.updateUserRed)

	// PUT
	mux.HandleFunc("PUT /api/users", api.updateUser)

	// DELETE
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", api.deleteChirp)

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server listening on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
