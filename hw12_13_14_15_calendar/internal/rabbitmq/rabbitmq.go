package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp" //nolint
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(dsn string) (*Client, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (c *Client) Publish(queue string, body []byte) error {
	return c.channel.Publish(
		"",    // exchange
		queue, // routing key (queue name)
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}
