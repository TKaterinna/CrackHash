package services

import (
	"context"
	"encoding/json"
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
	for _, task := range tasks {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Printf("Failed to marshal task %+v", task)
			return err
		}

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
		if err != nil {
			log.Printf("Failed to publish a message: %s", err)
			return err
		}

		log.Printf("SENT TASK %s", taskJSON)
		// TODO: а нужно ли эти задачи сохранять в базу?
	}

	return nil
}
