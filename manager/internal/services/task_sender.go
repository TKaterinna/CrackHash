package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
)

type TaskSender struct {
	client       *http.Client
	workersUrl   []string
	workersCount int64
}

func NewTaskSender(workersCount int64, workersPort []string) *TaskSender {
	var workersUrl []string

	for i := range workersCount {
		workersUrl = append(workersUrl, "http://worker"+strconv.Itoa(int(i))+workersPort[i]+"/internal/api/worker/hash/crack/task")
	}
	log.Println(workersUrl)

	return &TaskSender{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		workersUrl:   workersUrl,
		workersCount: workersCount,
	}
}

func (t *TaskSender) Send(tasks []*models.CrackTaskRequest) error {
	i := 0

	for _, task := range tasks {
		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Printf("Failed to marshal task %+v", task)
			return err
		}

		log.Println("SEND ", t.workersUrl[i])
		resp, err := t.client.Post(
			t.workersUrl[i],
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

		i += 1
		i = i % int(t.workersCount)
	}

	return nil
}
