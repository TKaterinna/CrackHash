package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQConnection struct {
	mu          sync.RWMutex
	Conn        *amqp.Connection
	Channel     *amqp.Channel
	URL         string
	reconnectCh chan struct{}
	closeCh     chan struct{}
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
		URL:     rabbitMQURL,
	}, nil
}

func (c *RMQConnection) NotifyReconnect() <-chan struct{} {
	return c.reconnectCh
}

func (c *RMQConnection) NotifyClose() <-chan struct{} {
	return c.closeCh
}

func (c *RMQConnection) SetupTopology() error {
	ch := c.GetChannel()
	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("channel not ready, waiting...")
	}
	err := ch.ExchangeDeclare(
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

	_, err = ch.QueueDeclare(
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

	err = ch.QueueBind(
		"task.queue",
		"task",
		"manager_worker",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
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

	err = ch.QueueBind(
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

func (c *RMQConnection) StartRecoveryWatcher(onReconnect func() error) {
	go func() {
		for {
			connNotify := c.Conn.NotifyClose(make(chan *amqp.Error, 1))

			c.mu.RLock()
			ch := c.Channel
			c.mu.RUnlock()

			var chanNotify chan *amqp.Error
			if ch != nil && !ch.IsClosed() {
				chanNotify = ch.NotifyClose(make(chan *amqp.Error, 1))
			}

			select {
			case err, ok := <-connNotify:
				if !ok {
					log.Println("Recovery watcher: connection permanently closed")
					close(c.closeCh)
					return
				}
				log.Printf("RabbitMQ connection lost: %v. Reconnecting...", err)
			case err, ok := <-chanNotify:
				if !ok {
					continue
				}
				log.Printf("RabbitMQ channel lost: %v. Reconnecting...", err)
			}

			var newConn *amqp.Connection
			var newCh *amqp.Channel
			var dialErr error

			for attempt := 1; attempt <= 10; attempt++ {
				log.Printf("Reconnect attempt %d/10...", attempt)

				newConn, dialErr = amqp.Dial(c.URL)
				if dialErr == nil {
					newCh, dialErr = newConn.Channel()
					if dialErr == nil {
						break
					}
					newConn.Close()
				}
				log.Printf("Reconnect failed: %v. Waiting 2s...", dialErr)
				time.Sleep(2 * time.Second)
			}

			if dialErr != nil {
				log.Fatalf("Failed to reconnect to RabbitMQ after 10 attempts: %v", dialErr)
			}

			c.mu.Lock()
			oldConn, oldCh := c.Conn, c.Channel
			c.Conn, c.Channel = newConn, newCh
			c.mu.Unlock()

			if oldCh != nil {
				_ = oldCh.Close()
			}
			if oldConn != nil {
				_ = oldConn.Close()
			}

			if err := c.SetupTopology(); err != nil {
				log.Printf("Failed to setup topology after reconnect: %v", err)
				continue
			}

			if onReconnect != nil {
				if err := onReconnect(); err != nil {
					log.Printf("Failed to re-subscribe listeners: %v", err)
				} else {
					log.Println("Listeners re-subscribed successfully")
				}
			}

			select {
			case c.reconnectCh <- struct{}{}:
			default:
			}

			log.Println("RabbitMQ reconnected and topology restored")
		}
	}()
}

func (c *RMQConnection) GetChannel() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Channel
}

func (c *RMQConnection) RecreateChannel() error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.Channel = ch
	c.mu.Unlock()
	return nil
}
