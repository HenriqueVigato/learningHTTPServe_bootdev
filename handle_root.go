package main

import "net/http"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./files/")).ServeHTTP(w, r)
}
