package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()
	return declareExchange(ch);
}

func  (e *Emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()

	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Pusing to channel...")

	err = ch.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}