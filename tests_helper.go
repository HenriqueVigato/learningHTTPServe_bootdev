package main

import (
	"context"
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
	testConfig.db.CreateUser(context.Background(), "user@test.test")

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM chirps")
		db.Close()
	})

	return &testConfig, nil
}
