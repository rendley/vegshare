package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	operationsModels "github.com/rendley/vegshare/backend/internal/operations/models"
	operationsRepository "github.com/rendley/vegshare/backend/internal/operations/repository"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/database"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
)

func main() {
	cfg := config.Load()

	// --- Инициализация зависимостей ---
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	client, err := rabbitmq.New(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer client.Close()

	opsRepo := operationsRepository.NewRepository(db)

	// --- Запуск консьюмера ---
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

			var opLog operationsModels.OperationLog
			if err := json.Unmarshal(d.Body, &opLog); err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				d.Nack(false, false) // Сообщение не может быть обработано, отбрасываем
				continue
			}

			// 1. Обновляем статус на "in_progress"
			log.Printf("Updating status to 'in_progress' for operation %s", opLog.ID)
			if err := opsRepo.UpdateOperationLogStatus(context.Background(), opLog.ID, "in_progress"); err != nil {
				log.Printf("Error updating status to in_progress: %s", err)
				d.Nack(false, true) // Ошибка, возможно временная, возвращаем в очередь
				continue
			}

			// 2. Имитируем работу
			log.Printf("Simulating work for 30 seconds...")
			time.Sleep(30 * time.Second)

			// 3. Обновляем статус на "completed"
			log.Printf("Updating status to 'completed' for operation %s", opLog.ID)
			if err := opsRepo.UpdateOperationLogStatus(context.Background(), opLog.ID, "completed"); err != nil {
				log.Printf("Error updating status to completed: %s", err)
				d.Nack(false, true) // Ошибка, возможно временная, возвращаем в очередь
				continue
			}

			log.Printf("Done processing message for operation %s.", opLog.ID)
			d.Ack(false) // Подтверждаем, что сообщение успешно обработано.
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}