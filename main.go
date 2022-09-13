package main

import (
	"log"

	"github.com/Lucas32-dev/monodirectional-chat-rabbitMQ/receiver"
	"github.com/Lucas32-dev/monodirectional-chat-rabbitMQ/sender"
)

// AMQP Server url
const serverURL = "amqp://guest:guest@localhost:5672/"

func main() {
	finished := make(chan struct{}, 1)

	// create sender
	s, err := sender.New(serverURL)

	// handle err
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	// create receiver
	r, err := receiver.New(serverURL)

	// handle err
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	msgs, err := r.ReadMessage("init")

	// handle err
	if err != nil {
		log.Fatalln(err)
	}

	// log messages
	go func() {
		for d := range msgs {
			log.Printf("\nmessage received: %s", d.Body)
			finished <- struct{}{}
		}
	}()

	// send message to queue 'init'
	err = s.SendMessage("init", "Hello world")

	// handle err
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Message sent with success!")

	// wait finish
	<-finished
	log.Println("Application finished")
}
