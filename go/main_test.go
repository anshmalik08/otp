package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"firebase.google.com/go"
	"google.golang.org/api/option"
)

func TestVerifyIDTokenHandler(t *testing.T) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("path/to/your/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		t.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	req := struct {
		IDToken string `json:"idToken"`
	}{
		IDToken: "TEST_ID_TOKEN",
	}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest("POST", "/verifyIdToken", bytes.NewReader(body))
	w := httptest.NewRecorder()

	verifyIDTokenHandler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if result["message"] != "ID token verified" {
		t.Errorf("Expected 'ID token verified' message but got %s", result["message"])
	}
}
