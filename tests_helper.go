package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/joho/godotenv"
)

func getConnectionTestDB(t *testing.T) (*apiConfig, error) {
	DB_URL := "postgres://user:pass123@localhost:5433/chirpy?sslmode=disable"
	err := godotenv.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("erro ao carregar as variaveis de ambiente: %v", err))
	}

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)

	testConfig := apiConfig{
		db:        dbQueries,
		plataform: os.Getenv("PLATAFORM"),
		secrete:   os.Getenv("SECRETE"),
	}

	db.Exec("DELETE FROM refresh_tokens")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM chirps")

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

	user1Params := database.CreateUserParams{
		Email:          "user1@test.test",
		HashedPassword: passwordHash,
	}
	user1, err := testConfig.db.CreateUser(context.Background(), user1Params)
	if err != nil {
		t.Fatalf("erro ao criar o usuario: %v", err)
	}

	if err = createChirps(&testConfig, user); err != nil {
		t.Fatalf("%v", err)
	}

	if err = createChirps(&testConfig, user1); err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM refresh_tokens")
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
			Body:   body + strconv.Itoa(rand.IntN(39999)),
			UserID: user.ID,
		})
		if err != nil {
			return fmt.Errorf("erro o adicionar o chirp (%s) no banco de dados: %v", body, err)
		}
	}
	return nil
}
