package sender

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Create a sender and connect to RabbitMQ server
func New(serverURL string) (*Sender, error) {
	s := &Sender{}

	// connection dial
	conn, err := amqp.Dial(serverURL)

	// handle err
	if err != nil {
		return nil, err
	}

	// save connection
	s.conn = conn

	// open a socket channel
	s.ch, err = conn.Channel()

	// handle err
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close AMQP connection
func (s *Sender) Close() {
	s.conn.Close()
}

// Send a string message to a queue
func (s *Sender) SendMessage(queueName string, message string) error {
	// Declare queue
	q, err := s.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	// handle err on declaring queue
	if err != nil {
		return err
	}

	// get context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish message
	err = s.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	// handle err on publishing
	if err != nil {
		return err
	}

	return nil
}
