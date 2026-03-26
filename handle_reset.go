package main

import (
	"log"
	"net/http"
	"strings"
)

func (api *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(api.plataform, "dev") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	api.fileserverHits.Store(0)
	err := api.db.DeleteUsers(r.Context())
	if err != nil {
		log.Fatalf("erro o deletar os usuarios %v", err)
	}
}
