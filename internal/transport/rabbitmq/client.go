package rabbitmqClient

import (
	"context"
	"encoding/json"
	"fmt"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"github.com/streadway/amqp"
)

type Client struct {
	conn *amqp.Connection
}

func NewClient(addr string) (*Client, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Failed to connect to RabbitMQ", err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) SendLogRequest(_ context.Context, req audit.LogItem) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to open a channel", err)
	}
	defer ch.Close()

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to marshal LogItem", err)
	}

	q, err := ch.QueueDeclare(
		"logs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to declare a queue", err)
	}

	if err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	); err != nil {
		return fmt.Errorf("%s: %w", "Failed to publish a message", err)
	}

	return nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}
