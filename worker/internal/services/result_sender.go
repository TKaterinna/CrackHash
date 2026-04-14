package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ResultSender struct {
	rabbit_conn *RMQConnection
}

func NewResultSender(rabbit_conn *RMQConnection) *ResultSender {
	return &ResultSender{
		rabbit_conn: rabbit_conn,
	}
}

func (r *ResultSender) Send(res *models.CrackTaskResult) error {
	ch := r.rabbit_conn.GetChannel()
	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("channel unavailable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal result %+v", res)
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"manager_worker",
		"result",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         resJSON,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
		return err
	}

	log.Printf("SENT RESULT %s", resJSON)

	return nil
}
