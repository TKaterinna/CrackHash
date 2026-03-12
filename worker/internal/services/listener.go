package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
)

type Listener struct {
	rabbit_conn *RMQConnection
	service     *CalcService
}

func NewCalcListener(rabbit_conn *RMQConnection, service *CalcService) *Listener {
	return &Listener{
		rabbit_conn: rabbit_conn,
		service:     service,
	}
}

func (l *Listener) Listen(ctx context.Context) {
	msgs, err := l.rabbit_conn.Channel.Consume(
		"task.queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Listener shutting down")
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("Messages channel closed")
				return
			}

			var req models.CrackTaskRequest
			log.Printf("READ %s", d.Body)
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Bad message: %v", err)
				d.Nack(false, false) // TODO: (что это?) Отбросить битое сообщение
				continue
			}

			if err := l.service.Save(&req); err != nil {
				log.Printf("Save request failed: %v", err)
				d.Nack(false, false) // TODO: (что это?) Отбросить битое сообщение
				continue
			}

			//ok
			d.Ack(false)
			log.Printf("Get task %s", req.TaskId)
		}
	}
}
