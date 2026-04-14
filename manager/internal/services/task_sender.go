package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskSender struct {
	rabbit_conn *RMQConnection
}

func NewTaskSender(rabbit_conn *RMQConnection) *TaskSender {
	return &TaskSender{
		rabbit_conn: rabbit_conn,
	}
}

func (t *TaskSender) Send(tasks []*models.CrackTaskRequest) error {
	var failed int

	for _, task := range tasks {
		if err := t.SendWithRetry(task, 6); err != nil {
			log.Printf("Failed to send task %s after retries: %v", task.TaskId, err)
			failed++
			continue
		}
		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()

		// taskJSON, err := json.Marshal(task)
		// if err != nil {
		// 	log.Printf("Failed to marshal task %s: %v", task.TaskId, err)
		// 	failed++
		// 	cancel()
		// 	continue
		// }

		// err = t.rabbit_conn.Channel.PublishWithContext(
		// 	ctx,
		// 	"manager_worker",
		// 	"task",
		// 	false,
		// 	false,
		// 	amqp.Publishing{
		// 		ContentType:  "application/json",
		// 		DeliveryMode: amqp.Persistent,
		// 		Body:         taskJSON,
		// 	},
		// )
		// if err != nil {
		// 	log.Printf("Failed to publish task %s: %v", task.TaskId, err)
		// 	failed++
		// 	continue
		// }

		log.Printf("SENT TASK %s", task.TaskId)
	}

	if failed > 0 {
		return fmt.Errorf("failed to send %d out of %d tasks", failed, len(tasks))
	}

	return nil
}

func (t *TaskSender) SendWithRetry(task *models.CrackTaskRequest, maxRetries int) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = t.rabbit_conn.Channel.PublishWithContext(
			ctx,
			"manager_worker",
			"task",
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         taskJSON,
			},
		)
		cancel()

		if err == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<attempt) // 0.5s, 1s, 2s, 4s...
			log.Printf("Publish attempt %d failed for task %s: %v. Retrying in %v...",
				attempt+1, task.TaskId, err, delay)
			time.Sleep(delay)
			continue
		}
	}

	return fmt.Errorf("publish failed after %d attempts: %w", maxRetries, err)
}
