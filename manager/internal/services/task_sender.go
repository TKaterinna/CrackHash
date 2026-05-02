package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
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

func (t *TaskSender) Send(tasks []*models.CrackTaskRequest) []uuid.UUID {
	var seccess []uuid.UUID

	for _, task := range tasks {
		if err := t.SendSingle(task); err != nil {
			log.Printf("Failed to send task %s after retries: %v", task.TaskId, err)
			continue
		}

		seccess = append(seccess, task.TaskId)

		log.Printf("SENT TASK %s", task.TaskId)
	}

	return seccess
}

func (t *TaskSender) SendSingle(task *models.CrackTaskRequest) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	ch := t.rabbit_conn.GetChannel()
	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("channel unavailable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
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
	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
		return err
	}

	log.Printf("SENT RESULT %s", taskJSON)

	return nil
}
