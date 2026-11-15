package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"eci-backend/internal/auth/model"
	"eci-backend/internal/platform/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(cfg.MongoURI),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(cfg.MongoDatabase).Collection("users")

	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD are required")
	}

	var existing model.User

	err = collection.FindOne(ctx, map[string]interface{}{
		"email": email,
	}).Decode(&existing)

	if err == nil {
		log.Println("Admin user already exists")
		return
	}

	if err != mongo.ErrNoDocuments {
		log.Fatal(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	user := model.User{
		Name:      "Administrator",
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := collection.InsertOne(ctx, user); err != nil {
		log.Fatal(err)
	}

	log.Println("Admin user seeded successfully")
}
