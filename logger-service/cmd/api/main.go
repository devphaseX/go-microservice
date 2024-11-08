package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	//connect to mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	//create a context in order to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	//close connection

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(mongoClient),
	}

	app.serve()
}

func (app *Config) serve() {
	routes := app.routes()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes,
	}

	fmt.Println(fmt.Sprintf("server listenining on port: %s", webPort))
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	opts := options.Client().ApplyURI(mongoURL)
	opts.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(opts)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return c, nil
}
