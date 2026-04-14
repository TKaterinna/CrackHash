package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Listener: shutting down")
				return
			default:
			}

			ch := l.rabbit_conn.GetChannel()
			if ch == nil || ch.IsClosed() {
				log.Println("Listener: channel not ready, waiting...")
				time.Sleep(2 * time.Second)
				continue
			}

			msgs, err := ch.Consume(
				"task.queue",
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("Listener: failed to consume: %v. Retrying...", err)
				time.Sleep(3 * time.Second)
				continue
			}

			log.Println("Listener: subscribed to task.queue")

			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-msgs:
					if !ok {
						log.Println("Listener: messages channel closed, will re-subscribe")
						break
					}

					var req models.CrackTaskRequest
					log.Printf("READ %s", d.Body)
					if err := json.Unmarshal(d.Body, &req); err != nil {
						log.Printf("Listener: bad message: %v", err)
						d.Nack(false, false) // сообщение отбрасывается: только оно + не возвращается в очередь
						continue
					}

					if err := l.service.Save(&req); err != nil {
						log.Printf("Save request failed: %v", err)
						d.Nack(false, false)
						continue
					}

					if err := d.Ack(false); err != nil {
						log.Printf("Listener: ack failed: %v", err)
					}
					log.Printf("Listener: processed task %s", req.TaskId)
				}
			}
		}
	}()
}
