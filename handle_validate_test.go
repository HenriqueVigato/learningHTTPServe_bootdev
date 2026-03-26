package main

import (
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {
	forbiddenWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cases := []string{
		"I heard that kerfuffle was quite bad",
		"What a KERFUFFLE that was yesterday",
		"He is such a Sharbert sometimes",
		"The stars are made of FORNAX and sharbert",
		"I had something interesting for breakfast",
		"The weather is nice today",
		"I love coding in Go",
		"kerfuffle",
		"KERFUFFLE SHARBERT FORNAX",
		"kerfuffles",
		"I like kerfuffle, sharbert, fornax",
	}

	for i, str := range cases {
		cases[i] = cleanInput(str, forbiddenWords)
	}

	if !strings.Contains(cases[0], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[0])
	}
	if !strings.Contains(cases[1], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[1])
	}
	if !strings.Contains(cases[2], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[2])
	}
	if !strings.Contains(cases[3], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[3])
	}
	if !strings.Contains(cases[7], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[7])
	}
	if !strings.Contains(cases[8], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[8])
	}
	if strings.Contains(cases[9], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[9])
	}
	if !strings.Contains(cases[10], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[10])
	}
}
