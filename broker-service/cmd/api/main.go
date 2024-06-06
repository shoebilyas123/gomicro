package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


type Config struct {
	Rabbit *amqp.Connection
}

const brokerAddr = "80"

func main(){
	// try to connect to rabbitMQ
	rabbitConn, err := connectToRabbitMQ()

	if err != nil {
		log.Println(err);
		os.Exit(1);
	}
	
	defer rabbitConn.Close()
	app := Config{
		Rabbit: rabbitConn,
	};

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",brokerAddr),
		Handler: app.routes(),
	}

	log.Printf("Server listening on PORT:%s\n",brokerAddr);
	err = srv.ListenAndServe();

	if err != nil {
		log.Fatal(err);
	}
	
}


func connectToRabbitMQ() (*amqp.Connection, error) {
	var count int64
	var backoff = 1*time.Second

	var connection *amqp.Connection
	
	// Don't continue until rabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672");

		if err != nil {
			log.Println("RabbitMQ is not yet ready");
			count++;
		}	else {
			log.Println("Connected to RabbitMQ!");
			connection = c;
			break;
		}

		if count > 5 {
			log.Println(err);
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(backoff),2)) * time.Second;
		log.Printf("Retrying connection in %s seconds\n", backoff);
		time.Sleep(backoff)
		continue;
	}

	return connection, nil;
}