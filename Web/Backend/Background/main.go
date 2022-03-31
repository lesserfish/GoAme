package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
	ame "github.com/lesserfish/GoAme/Ame"
	"github.com/streadway/amqp"
)

var (
	workercount       uint
	AMQPURI           string
	AMQPIP            uint
	REDISURI          string
	REDISIP           uint
	REDISPROC         string
	queuename         string
	IDshift           uint
	configuration     string
	StorageDirectory  string
	DownloadDirectory string
	PersistenceTime   float64
	CleanTime         float64
)

func main() {

	flag.UintVar(&workercount, "n", 64, "Specify the quantity of workers to be used.")
	flag.StringVar(&AMQPURI, "amqpurl", "amqp://localhost", "Address of RabbitMQ server")
	flag.StringVar(&REDISURI, "redisurl", "localhost", "Address of RabbitMQ server")
	flag.UintVar(&AMQPIP, "amqpp", 5672, "port of the RabbitMQ server")
	flag.UintVar(&REDISIP, "redisp", 6379, "port of the Redis server")
	flag.StringVar(&REDISPROC, "redisproc", "tcp", "Redis protocol. 'tcp' or 'udp'")
	flag.StringVar(&queuename, "queue", "ame", "Name of the queue to be used!")
	flag.UintVar(&IDshift, "idshift", 0, "Value of the starting ID of the workers")
	flag.StringVar(&configuration, "c", "", "Configuration file for Ame")
	flag.StringVar(&StorageDirectory, "storage", "/tmp", "Directory for storage of temporary files.")
	flag.StringVar(&DownloadDirectory, "download", "/tmp", "Directory for storage of download files.")
	flag.Float64Var(&PersistenceTime, "persistence", 30, "Time zip files will stick around in minutes")
	flag.Float64Var(&CleanTime, "cleanperiod", 10, "Period for files to be deleted")
	flag.Parse()

	if len(configuration) == 0 {
		log.Println("You need to specify a configuration file for Ame.")
		flag.PrintDefaults()
		return
	}

	// Initialize AMQP

	fulladdr := AMQPURI + ":" + strconv.Itoa(int(AMQPIP))
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
	err = channel.Qos(1, 0, false)

	if err != nil {
		log.Fatalln(err)
	}

	// Initialize Redis

	redisPool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(REDISPROC, REDISURI+":"+strconv.Itoa(int(REDISIP)))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	redisClient := redisPool.Get()
	defer redisClient.Close()

	// Initialize Ame

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

	// Create workers

	for id := uint(1); id <= workercount; id++ {

		newworker := Worker{id + IDshift,
			channel,
			queue.Name,
			redisClient,
			ameinstance}

		go newworker.Work()
	}

	// Start Cleaners

	go CleanTasker()

	forever := make(chan bool)
	<-forever
}
