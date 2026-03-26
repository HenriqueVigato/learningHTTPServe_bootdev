package main

import "testing"

func TestAddUser(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Logf("erro ao se conenctar ao bando de dados %v", err)
	}
}
