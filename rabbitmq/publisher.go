package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func PublishMessage(ch *amqp091.Channel, q amqp091.Queue, message []byte) error {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return fmt.Errorf("Failed to publish a message to RabbitMQ: %s", err)
	}
	return nil
}
