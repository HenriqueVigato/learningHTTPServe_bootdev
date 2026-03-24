package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (api *apiConfig) middlewareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (api *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", api.fileserverHits.Load())))
}

func (api *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	api.fileserverHits.Store(0)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./files/")).ServeHTTP(w, r)
}

func handleLogo(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./files/assets/")).ServeHTTP(w, r)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	var apiConf apiConfig
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiConf.middlewareMetricsInt(http.HandlerFunc(handleRoot))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", apiConf.middlewareMetricsInt(http.HandlerFunc(handleLogo))))
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /api/metrics", apiConf.metrics)
	mux.HandleFunc("POST /api/reset", apiConf.reset)

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server listening on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
