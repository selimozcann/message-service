package rabbitmq

import "github.com/rabbitmq/amqp091-go"

func ConsumeMessage(ch *amqp091.Channel, q amqp091.Queue) (<-chan amqp091.Delivery, error) {
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return msgs, err
}
