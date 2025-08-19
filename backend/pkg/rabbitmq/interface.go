package rabbitmq

import "github.com/streadway/amqp"

// ClientInterface defines the contract for a RabbitMQ client.
// This allows for mocking in tests.
type ClientInterface interface {
	Publish(queueName, body string) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Close()
}
