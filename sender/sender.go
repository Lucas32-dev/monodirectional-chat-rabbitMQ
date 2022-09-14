package sender

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	conn  *amqp.Connection
	chats []Chat
}

type Chat struct {
	name string
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

	// handle err
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close AMQP connection
func (s *Sender) Close() {
	// close chat connections
	for _, c := range s.chats {
		c.ch.Close()
	}

	// close connection
	s.conn.Close()
}
