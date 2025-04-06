package main

import (
	"blogService/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8081"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	log.Println("mongoclient:")

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	log.Println("1:")

	app := Config{
		Models: data.New(client),
	}
	log.Println("2:")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Println("3:")

	fmt.Printf("hello")
	err = srv.ListenAndServe()
	log.Println("errr")

	if err != nil {
		fmt.Printf("err")

		log.Println("Error starting server:", err)
		log.Panic()
	}
	log.Println("lets go")

}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	fmt.Println("Mongo connection")
	log.Println("Mongo:")

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)

		return nil, err
	}
	log.Println("Connected to mongo!")

	return c, nil
}
