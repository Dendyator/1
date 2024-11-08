package rabbitmq_test

import (
	"testing"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/rabbitmq" //nolint
	"github.com/stretchr/testify/assert"
)

func TestRabbitMQConnection(t *testing.T) {
	client, err := rabbitmq.New("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = client.DeclareQueue("test_queue")
	assert.NoError(t, err)

	err = client.Publish("test_queue", []byte("test message"))
	assert.NoError(t, err)

	client.Close()
}
