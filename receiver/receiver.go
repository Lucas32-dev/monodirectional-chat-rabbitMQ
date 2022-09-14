package receiver

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Receiver struct {
	conn        *amqp.Connection
	chatReaders []ChatReader
}

type ChatReader struct {
	name string
	ch   *amqp.Channel
	q    amqp.Queue
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

// Read messages
func (c ChatReader) ReadMessages() (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		c.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}

// Creates a chat receiver
func (r *Receiver) ConnectChat(name string, queueName string) (ChatReader, error) {
	// create socket channel
	ch, err := r.conn.Channel()

	// handle err
	if err != nil {
		return ChatReader{}, err
	}

	// create queue
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	// handle err
	if err != nil {
		return ChatReader{}, err
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,
		"",
		name,
		false,
		nil,
	)

	// handle err
	if err != nil {
		return ChatReader{}, err
	}

	// Create chat reader
	chatReader := ChatReader{
		name: name,
		ch:   ch,
		q:    q,
	}

	// save chat reader
	r.chatReaders = append(r.chatReaders, chatReader)

	return chatReader, err
}
