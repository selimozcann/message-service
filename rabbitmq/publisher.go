package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

func PublishMessage(ch *amqp091.Channel, q amqp091.Queue, message []byte) error {
	return ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
}
