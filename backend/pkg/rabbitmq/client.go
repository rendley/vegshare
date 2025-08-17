
package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Client is a RabbitMQ client.
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// New creates a new RabbitMQ client.
func New(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Client{conn: conn, ch: ch}, nil
}

// Close closes the RabbitMQ connection and channel.
func (c *Client) Close() {
	c.ch.Close()
	c.conn.Close()
}

// Publish publishes a message to a queue.
func (c *Client) Publish(queueName, body string) error {
	q, err := c.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = c.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

// Consume consumes messages from a queue.
func (c *Client) Consume(queueName string) (<-chan amqp.Delivery, error) {
	q, err := c.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return msgs, nil
}
