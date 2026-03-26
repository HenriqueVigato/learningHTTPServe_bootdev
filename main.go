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
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./files/")).ServeHTTP(w, r)
}

func handleLogo(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./files/assets/")).ServeHTTP(w, r)
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

	apiConf := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}
	mux := http.NewServeMux()

	// GET
	mux.Handle("/app/", http.StripPrefix("/app", apiConf.middlewareMetricsInt(http.HandlerFunc(handleRoot))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", apiConf.middlewareMetricsInt(http.HandlerFunc(handleLogo))))
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /admin/metrics", apiConf.metrics)

	// POST
	mux.HandleFunc("POST /admin/reset", apiConf.reset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server listening on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
