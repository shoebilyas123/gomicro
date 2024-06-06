package main

import (
	"context"
	"fmt"
	"log"
	"logsvc/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	mongoURL = "mongodb://mongo:27017"
	// gRPCPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connec to Mongo

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second);
	defer cancel()

	mongoClient, err := connectToMongo(ctx);

	log.Println("Connected to mongo-client");

	if err != nil {
		log.Panic(err);
	}

	client = mongoClient

	defer func(){
		if err:=client.Disconnect(ctx); err != nil {
			log.Panic(err);
		}
	}()

	if err != nil {
		log.Panic(err);
	}
	app := Config{data.New(client)}

	// Register the RPC Server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Println("Cannot register RPC Server");
		
	}

	go app.rpcListen()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}

	log.Printf("Logger service running on PORT:%s\n",webPort)
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err);
	}
	
}

func (app *Config) rpcListen() error {
	log.Printf("Starting RPC server on PORT:%s\n",rpcPort)
	
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s",rpcPort))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			return err
		}
		go rpc.ServeConn(rpcConn);
	}
}

func connectToMongo(ctx context.Context) (*mongo.Client, error) {

	mongoOptions := options.Client().ApplyURI(mongoURL);
	mongoOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(ctx, mongoOptions);
	
	if err != nil {
		log.Println("Error connecting to MongoDB",err);
		return nil, err;
	}

	return client, nil;
}

