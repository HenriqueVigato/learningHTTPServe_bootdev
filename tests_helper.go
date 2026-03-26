package main

import (
	"database/sql"
	"testing"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
)

func getConnectionTestDB(t *testing.T) (*apiConfig, error) {
	DB_URL := "postgres://user:pass123@localhost:5433/chirpy?sslmode=disable"

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)

	testConfig := apiConfig{
		db:        dbQueries,
		plataform: "dev",
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Close()
	})

	return &testConfig, nil
}
