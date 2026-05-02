package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/TKaterinna/CrackHash/worker/config"
	"github.com/TKaterinna/CrackHash/worker/internal/repo"
	"github.com/TKaterinna/CrackHash/worker/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := config.NewConfig()

	calcRepo := repo.NewCalcRepo()
	rabbit_conn, err := services.RabbitMQConnect(config.RabbitMQURL)
	if err != nil {
		log.Panicf("Failed to connect RabbitMQ: %s", err)
	}

	err = rabbit_conn.SetupTopology()
	if err != nil {
		log.Panicf("Failed to setup topology: %s", err)
	}

	resultSender := services.NewResultSender(rabbit_conn)
	calcService := services.NewCalcService(calcRepo, resultSender, config.SleepMs)

	rabbit_conn.StartRecoveryWatcher(func() error {
		return nil
	})

	listener := services.NewCalcListener(rabbit_conn, calcService)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server starting on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		listener.Listen(context.Background())
	}()

	wg.Wait()
	log.Println("Stop worker")
}
