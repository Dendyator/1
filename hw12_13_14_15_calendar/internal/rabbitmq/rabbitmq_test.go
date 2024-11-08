package rabbitmq_test

import (
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRabbitMQConnection(t *testing.T) {
	var client *rabbitmq.Client
	var err error

	for i := 0; i < 5; i++ {
		client, err = rabbitmq.New("amqp://guest:guest@localhost:5672/")
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = client.DeclareQueue("test_queue")
	assert.NoError(t, err)

	err = client.Publish("test_queue", []byte("test message"))
	assert.NoError(t, err)

	client.Close()
}
