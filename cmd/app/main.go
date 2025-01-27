package main

import (
	"context"
	"github.com/aarsinh/auth-golang/middleware"
	"github.com/aarsinh/auth-golang/pkg/models/mongodb"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *mongodb.UserModel
}

func main() {
	// Initialize loggers for info and error
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	err := godotenv.Load()
	if err != nil {
		errLog.Printf("Could not load .env file")
	}

	middleware.InitMiddleware()

	// Open database connection
	mongoURI := "mongodb://localhost:27017/"
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	co := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, co)

	if err != nil {
		errLog.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("authentication").Collection("users")
	infoLog.Printf("Database connection established")

	app := &application{
		errorLog: errLog,
		infoLog:  infoLog,
		db:       &mongodb.UserModel{C: collection},
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
		ErrorLog:     errLog,
	}
	infoLog.Printf("Starting server on port 8008")
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
