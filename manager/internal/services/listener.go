package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
)

type Listener struct {
	rabbit_conn *RMQConnection
	service     *TaskService
}

func NewCalcListener(rabbit_conn *RMQConnection, service *TaskService) *Listener {
	return &Listener{
		rabbit_conn: rabbit_conn,
		service:     service,
	}
}

func (l *Listener) Listen(ctx context.Context) {
	msgs, err := l.rabbit_conn.Channel.Consume(
		"result.queue",
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

	go func() {
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

				var req models.CrackTaskResult
				log.Printf("READ %s", d.Body)
				if err := json.Unmarshal(d.Body, &req); err != nil {
					log.Printf("Bad message: %v", err)
					d.Nack(false, false)
					continue
				}

				if err := l.service.UpdateResult(&req); err != nil {
					log.Printf("Update result in db failed: %v", err) // вроде может сработать при дубликате, например дошло сообщение от выпавшего воркера, а другой уже досчитал эту таску
					d.Nack(false, false)
					continue
				}

				d.Ack(false)
				log.Printf("Result processed for task %s", req.TaskId)
			}
		}
	}()
}
