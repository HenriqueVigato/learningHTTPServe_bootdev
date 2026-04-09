package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
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
		secrete:   "z4Otc/8ibsFtsaCLbfPm8xvmWUuAsz05L9Qr4UNZQ9Nm0Rqk+5qWV46II5zOPM4XMMVLYd8s+jkPzSWe6qzT9w==",
	}
	passwordHash, err := auth.HashPassword("password")
	if err != nil {
		t.Fatalf("erro ao hashear a senha %v", err)
	}

	userParams := database.CreateUserParams{
		Email:          "user@test.test",
		HashedPassword: passwordHash,
	}
	user, err := testConfig.db.CreateUser(context.Background(), userParams)
	if err != nil {
		t.Fatalf("erro ao criar o usuario: %v", err)
	}

	if err = createChirps(&testConfig, user); err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM chirps")
		db.Close()
	})

	return &testConfig, nil
}

func createChirps(api *apiConfig, user database.User) error {
	chirps := []string{
		"Hello, testing the handler chirps",
		"I heard that kerfuffle was quite bad",
		"Did you hear about that sharbert incident?",
		"The fornax was mentioned in the report",
		"Just a normal chirp to fill the database",
		"Another chirp for testing purposes",
	}

	for _, body := range chirps {
		_, err := api.db.CreateChirp(context.Background(), database.CreateChirpParams{
			Body:   body,
			UserID: user.ID,
		})
		if err != nil {
			return fmt.Errorf("erro o adicionar o chirp (%s) no banco de dados: %v", body, err)
		}
	}
	return nil
}
