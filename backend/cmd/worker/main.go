
package main

import (
	"log"

	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
)

func main() {
	cfg := config.Load()

	client, err := rabbitmq.New(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer client.Close()

	msgs, err := client.Consume(cfg.RabbitMQ.Queues["actions"])
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Здесь в будущем будет бизнес-логика обработки.
			// Например, отправка уведомления в Телеграм.
			// Если эта логика вернет ошибку, нужно будет вызвать d.Nack()

			log.Printf("Done processing message.")
			d.Ack(false) // Подтверждаем, что сообщение успешно обработано.
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
