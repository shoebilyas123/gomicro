package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()

	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()

	if err != nil {
		return err;
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()

	if err != nil {
		return err;
	}
	defer ch.Close()

	q, err := declareQueue(ch);
	if err != nil {
		return err;
	}

	for _, s := range topics {
		err = ch.QueueBind(q.Name, s, "logs_topic", false, nil)

		if err != nil {
			return err;
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err;
	}

	forever := make(chan bool);

	go func(){
		for d := range msgs {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload);
			go handlePayload(payload)
		}
	}()
		
	fmt.Printf("Waiting for message on exchange queue...");
	<-forever

	return nil;
}

func handlePayload(payload Payload) {
	switch payload.Name {
		case "log","event":
			err := logEvent(payload);

			if err != nil {
				log.Println(err);
			}
			break;
		case "auth":
			// something on auth
			break;
		default:
			err := logEvent(payload);

			if err != nil {
				log.Println(err);
			}
	}
}

func logEvent(logData Payload) error {
	jsonData, _ := json.MarshalIndent(logData, "", "\t");

	reqBody := bytes.NewBuffer(jsonData);
	request, err := http.NewRequest("POST", "http://logsvc/log", reqBody);
	request.Header.Set("Content-Type","application/json");

	if err != nil {
		return err;
	}

	client := &http.Client{};
	response, err := client.Do(request);

	if err != nil {
		return err;
	}

	defer response.Body.Close();

		// Verify the response from svc and write the appropriate response to the calling client
	if response.StatusCode != http.StatusAccepted {
		return err;
	}
	return nil;
}