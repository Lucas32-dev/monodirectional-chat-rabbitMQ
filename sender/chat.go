package sender

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Creates a chat
func (s *Sender) CreateChat(name string) (Chat, error) {
	// create socket channel
	ch, err := s.conn.Channel()

	// handle err
	if err != nil {
		return Chat{}, err
	}

	// exange declare
	err = ch.ExchangeDeclare(
		name,     // name
		"fanout", // type
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	// handle err
	if err != nil {
		return Chat{}, err
	}

	chat := Chat{
		name: name,
		ch:   ch,
	}

	// save chat
	s.chats = append(s.chats, chat)

	return chat, nil
}

// Sends a string message to an exchange
func (c Chat) SendMessage(message string) error {
	// get ctx
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.ch.PublishWithContext(ctx,
		c.name, // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
