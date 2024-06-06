package main

import (
	"context"
	"log"
	"logsvc/data"
	"time"
)

type RPCServer struct {}

type RPCPayload struct {
	Data string
	Name string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs");

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})	
	if err != nil {
		log.Println("Error writing to Mongo")
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name;
	return nil;
}