package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App

func main() {
	ctx := context.Background()

	
	sa := option.WithCredentialsFile("otp-project-85e8a-firebase-adminsdk-q4n09-9893b87ec6.json") 
	var err error
	app, err = firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	http.HandleFunc("/verifyIdToken", verifyIDTokenHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func verifyIDTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	auth, err := app.Auth(ctx)
	if err != nil {
		http.Error(w, "Firebase Auth error", http.StatusInternalServerError)
		return
	}

	
	var req struct {
		IDToken string `json:"idToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	
	token, err := auth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token", http.StatusUnauthorized)
		return
	}

	
	resp := map[string]string{"message": "ID token verified", "uid": token.UID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
