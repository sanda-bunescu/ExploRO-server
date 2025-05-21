package initializers

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
)

var FirebaseApp *firebase.App

func FirebaseInitialization() {
	opt := option.WithCredentialsFile("/app/secrets/firebase-creds.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to connect to initialize firebase %v", err)
	}
	FirebaseApp = app
}
