package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	http.FileServer(http.Dir("./files/")).ServeHTTP(w, r)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", apiHandler{})

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server listening on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
