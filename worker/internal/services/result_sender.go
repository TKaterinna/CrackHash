package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
)

type ResultSender struct {
	client     *http.Client
	managerUrl string
}

func NewResultSender(managerPort string) *ResultSender {
	return &ResultSender{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		managerUrl: "http://manager" + managerPort + "/internal/api/manager/hash/crack/request",
	}
}

func (r *ResultSender) Send(res *models.CrackTaskResult) error {
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal result %+v", res)
		return err
	}

	req, err := http.NewRequest(
		http.MethodPatch,
		r.managerUrl,
		bytes.NewBuffer(resJSON),
	)
	if err != nil {
		log.Printf("Failed to send result %+v: %v", res, err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Received non-success status %d for result %+v", resp.StatusCode, res)
	}

	return nil
}
