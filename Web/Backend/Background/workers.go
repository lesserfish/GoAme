package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	ame "github.com/lesserfish/GoAme/Ame"
	"github.com/streadway/amqp"
)

type Worker struct {
	workerID  uint
	channel   *amqp.Channel
	queueName string
	AmeKanji  *ame.AmeKanji
}

type jstest struct {
	Name []string
	Age  int
}

func (worker Worker) Work() {
	logmsg := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] Starting work!"
	log.Println(logmsg)

	msgs, err := worker.channel.Consume(
		worker.queueName,
		strconv.Itoa(int(worker.workerID)),
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		errmsg := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] "
		errmsg += "Error: " + err.Error()
		log.Println(errmsg)
		return
	}

	for msg := range msgs {
		var tst jstest
		json.Unmarshal(msg.Body, &tst)
		fmt.Println("Yo! The value was " + strconv.Itoa(tst.Age))
		time.Sleep(3 * time.Second)
		msg.Ack(false)
	}

	text := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] " + "Exiting!"
	fmt.Println(text)

}
