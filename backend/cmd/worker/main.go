package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	catalogService "github.com/rendley/vegshare/backend/internal/catalog/service"
	operationsModels "github.com/rendley/vegshare/backend/internal/operations/models"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"

	catalogRepository "github.com/rendley/vegshare/backend/internal/catalog/repository"
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	farmService "github.com/rendley/vegshare/backend/internal/farm/service"
	operationsRepository "github.com/rendley/vegshare/backend/internal/operations/repository"
	plotRepository "github.com/rendley/vegshare/backend/internal/plot/repository"
	taskRepository "github.com/rendley/vegshare/backend/internal/task/repository"
	taskService "github.com/rendley/vegshare/backend/internal/task/service"
	unitcontentRepository "github.com/rendley/vegshare/backend/internal/unitcontent/repository"
	unitcontentService "github.com/rendley/vegshare/backend/internal/unitcontent/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/database"
	"github.com/rendley/vegshare/backend/pkg/logger"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
)

// PlantActionParams - структура для парсинга параметров операции 'plant'.
type PlantActionParams struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

func main() {
	cfg := config.Load()
	log := logger.New()

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

	// Инициализируем все репозитории и сервисы, которые нужны воркеру-оркестратору
	opsRepo := operationsRepository.NewRepository(db)
	taskRepo := taskRepository.NewRepository(db)
	unitContentRepo := unitcontentRepository.NewRepository(db)
	plotRepo := plotRepository.NewRepository(db)
	farmRepo := farmRepository.NewRepository(db)
	catalogRepo := catalogRepository.NewRepository(db)

	unitContentSvc := unitcontentService.NewService(unitContentRepo)
	farmSvc := farmService.NewFarmService(farmRepo)
	plotSvc := plotService.NewService(plotRepo, farmSvc)
	catalogSvc := catalogService.NewService(catalogRepo)
	taskSvc := taskService.NewService(db, taskRepo, opsRepo, unitContentSvc, plotSvc, catalogSvc)

	// --- Запуск консьюмера ---
	queueName := cfg.RabbitMQ.Queues["actions"]
	msgs, err := client.Consume(queueName)
	if err != nil {
		log.Fatalf("Failed to register a consumer for queue '%s': %v", queueName, err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Infof("Received a message: %s", d.Body)

			var opLog operationsModels.OperationLog
			if err := json.Unmarshal(d.Body, &opLog); err != nil {
				log.Errorf("Error unmarshalling message: %s", err)
				d.Nack(false, false)
				continue
			}

			ctx := context.Background()

			// 1. Обновляем статус операции на 'processing'
			log.Infof("Updating operation %s to 'processing'", opLog.ID)
			if err := opsRepo.UpdateOperationLogStatus(ctx, opLog.ID, "processing"); err != nil {
				log.Errorf("Error updating operation status: %s", err)
				d.Nack(false, true)
				continue
			}

			// 2. Формируем заголовок задачи
			title, err := buildTaskTitle(ctx, &opLog, plotSvc, catalogSvc)
			if err != nil {
				log.Errorf("Error building task title: %s", err)
				d.Nack(false, true)
				continue
			}

			// 3. Создаем задачу
			log.Infof("Creating task for operation %s with title: %s", opLog.ID, title)
			_, err = taskSvc.CreateTask(ctx, opLog.ID, title, string(opLog.Parameters))
			if err != nil {
				log.Errorf("Error creating task: %s", err)
				d.Nack(false, true)
				continue
			}

			// 4. TODO: Отправить уведомление в Telegram

			log.Infof("Successfully created task for operation %s.", opLog.ID)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func buildTaskTitle(ctx context.Context, operation *operationsModels.OperationLog, plotSvc plotService.Service, catalogSvc catalogService.Service) (string, error) {
	var plotName string
	if operation.UnitType == "plot" {
		plot, err := plotSvc.GetPlotByID(ctx, operation.UnitID)
		if err != nil {
			return "", err
		}
		plotName = plot.Name
	} else {
		plotName = operation.UnitID.String()
	}

	if operation.ActionType == "plant" {
		var params PlantActionParams
		if err := json.Unmarshal(operation.Parameters, &params); err != nil {
			return "", fmt.Errorf("ошибка парсинга параметров для операции plant: %w", err)
		}
		item, err := catalogSvc.GetItemByID(ctx, params.ItemID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Посадить '%s' (x%d) на грядке '%s'", item.Name, params.Quantity, plotName), nil
	}

	return fmt.Sprintf("Выполнить '%s' на '%s'", operation.ActionType, plotName), nil
}