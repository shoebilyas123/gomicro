package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)


func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		"logs_topic",		//name?
		"topic",				// what topic are we reading?
		true,				// Is this exchange durable?
		false,				// Do you get rid of this channel when you are done?
		false,				// false because we are using this across microservices
		false,				// no-wait?
		nil,				// arguments??
	)
}


func declareQueue(channel *amqp.Channel) (amqp.Queue,error) {
	return channel.QueueDeclare(
		"myqueue",	//name?
		false,		// durable?
		false,		// delete when unused?
		true,		// exclusive?
		false,		// no-wait
		nil,
	)
}