package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type jstest struct {
	Name []string
	Age  int
}

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
		"ame",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err)
	}

	text := jstest{[]string{"a", "b", "c"}, 121}
	byteinfo, err := json.Marshal(text)

	fmt.Println(string(byteinfo))

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteinfo,
	})

	if err != nil {
		log.Println(err)
	}
}
