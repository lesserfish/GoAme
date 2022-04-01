package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
	ame "github.com/lesserfish/GoAme/Ame"
	"github.com/streadway/amqp"
)

func old_main() {
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

	input_file := "/home/lesserfish/Documents/Code/GoAme/Resources/input.json"
	input_content, err := ioutil.ReadFile(input_file)

	if err != nil {
		log.Fatalln(err)
	}
	var input ame.Input

	err = json.Unmarshal(input_content, &input)

	if err != nil {
		log.Fatalln(err)
	}

	info := Message{}
	info.Input = input
	info.UUID, _ = uuid.NewUUID()

	byteinfo, err := json.Marshal(info)

	fmt.Println(string(byteinfo))

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteinfo,
	})

	if err != nil {
		log.Println(err)
	}
}
