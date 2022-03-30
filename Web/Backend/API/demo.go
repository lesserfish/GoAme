package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		log.Println(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Println(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"first",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err)
	}

	text := "hello world!"
	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(text),
	})

	if err != nil {
		log.Println(err)
	}
}
