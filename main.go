package main

import (
	"fmt"
	"log"
	"net/http"
)

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
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.HandlerFunc(handleRoot)))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.HandlerFunc(handleLogo)))
	mux.HandleFunc("/healthz/", handleHealth)

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server listening on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
