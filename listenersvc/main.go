package main

import (
	"listenersvc/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitMQ
	rabbitConn, err := connectToRabbitMQ()

	if err != nil {
		log.Println(err);
		os.Exit(1);
	}

	defer rabbitConn.Close()
	// Start listening for messages
	// To create consumer

	consumer, err := event.NewConsumer(rabbitConn)

	if err != nil {
		
		log.Panic(err);
	}
	// Watch the queue and consume events from topics
	err = consumer.Listen([]string{"logs.INFO","logs.Error","logs.Warning"});

	if err != nil {
		log.Println(err);
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