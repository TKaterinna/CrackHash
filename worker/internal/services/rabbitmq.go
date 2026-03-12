package services

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func RabbitMQConnect(rabbitMQURL string) (*RMQConnection, error) {
	var conn *amqp.Connection
	var err error
	var maxRetries = 10
	var retryDelay = 2 * time.Second

	for i := range maxRetries {
		log.Printf("Attempting to connect to RabbitMQ (attempt %d/%d): %s", i+1, maxRetries, rabbitMQURL)

		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			log.Println("Successfully connected to RabbitMQ")
			break
		}

		log.Printf("Connection failed: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to make channel: %w", err)
	}

	return &RMQConnection{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (c *RMQConnection) SetupTopology() error {
	err := c.Channel.ExchangeDeclare(
		"manager_worker",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.Channel.QueueDeclare(
		"task.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.Channel.QueueBind(
		"task.queue",
		"task",
		"manager_worker",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.Channel.QueueDeclare(
		"result.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.Channel.QueueBind(
		"result.queue",
		"result",
		"manager_worker",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
