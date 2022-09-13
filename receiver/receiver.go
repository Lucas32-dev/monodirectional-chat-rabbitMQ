package receiver

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Receiver struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Create a receiver and connect to RabbitMQ server
func New(serverURL string) (*Receiver, error) {
	r := &Receiver{}

	// connection dial
	conn, err := amqp.Dial(serverURL)

	// handle err
	if err != nil {
		return nil, err
	}

	// save connection
	r.conn = conn

	// open a socket channel
	r.ch, err = conn.Channel()

	// handle err
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Close AMQP connection
func (s *Receiver) Close() {
	s.conn.Close()
}

// Read messages to a queue
func (s *Receiver) ReadMessage(queueName string) (<-chan amqp.Delivery, error) {
	return s.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}
