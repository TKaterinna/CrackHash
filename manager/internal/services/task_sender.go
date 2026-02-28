package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
)

type TaskSender struct {
	client    *http.Client
	workerUrl string
}

func NewTaskSender(workerPort string) *TaskSender {
	return &TaskSender{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		workerUrl: "http://worker" + workerPort + "/internal/api/worker/hash/crack/task",
	}
}

func (t *TaskSender) Send(tasks []*models.CrackTaskRequest) error {
	for _, task := range tasks {
		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Printf("Failed to marshal task %+v", task)
			return err
		}

		resp, err := t.client.Post(
			t.workerUrl,
			"application/json",
			bytes.NewBuffer(taskJSON),
		)
		if err != nil {
			log.Printf("Failed to send task %+v: %v", task, err)
			return err
		}

		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Received non-success status %d for task %+v", resp.StatusCode, task)
		}
	}

	return nil
}
