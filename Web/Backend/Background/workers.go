package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
		text := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] " + string(msg.Body)
		fmt.Println(text)
		time.Sleep(3 * time.Second)
		msg.Ack(false)
	}

}

func main() {
	var workercount uint
	var URI string
	var IP uint
	var queuename string
	var IDshift uint
	var configuration string

	flag.UintVar(&workercount, "n", 64, "Specify the quantity of workers to be used.")
	flag.StringVar(&URI, "url", "amqp://localhost", "Address of RabbitMQ server")
	flag.UintVar(&IP, "p", 5672, "port of the RabbitMQ server")
	flag.StringVar(&queuename, "queue", "ame", "Name of the queue to be used!")
	flag.UintVar(&IDshift, "idshift", 0, "Value of the starting ID of the workers")
	flag.StringVar(&configuration, "c", "", "Configuration file for Ame")

	flag.Parse()

	if len(configuration) == 0 {
		log.Println("You need to specify a configuration file for Ame.")
		flag.PrintDefaults()
		return
	}

	fulladdr := URI + ":" + strconv.Itoa(int(IP))
	connection, err := amqp.Dial(fulladdr)

	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queuename,
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Fatalln(err)
	}

	config_content, err := ioutil.ReadFile(configuration)

	if err != nil {
		log.Fatalln(err)
	}

	var config ame.Configuration
	json.Unmarshal(config_content, &config)

	ameinstance, err := ame.Initialize(config)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Ame initialized!")

	for id := uint(1); id <= workercount; id++ {
		newworker := Worker{id + IDshift, channel, queue.Name, ameinstance}
		go newworker.Work()
	}

	forever := make(chan bool)
	<-forever
}
