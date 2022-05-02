package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
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
	DownloadDirectory string
	StorageDirectory  string
	PersistenceTime   float64
	CleanTime         float64
	ctx               context.Context
)

func main() {

	flag.UintVar(&workercount, "n", 64, "Specify the quantity of workers to be used.")
	flag.StringVar(&AMQPURI, "amqp", "amqp://localhost", "Address of RabbitMQ server")
	flag.StringVar(&REDISURI, "redis", "localhost", "Address of RabbitMQ server")
	flag.UintVar(&AMQPIP, "amqpp", 5672, "port of the RabbitMQ server")
	flag.UintVar(&REDISIP, "redisp", 6379, "port of the Redis server")
	flag.StringVar(&REDISPROC, "redisproc", "tcp", "Redis protocol. 'tcp' or 'udp'")
	flag.StringVar(&queuename, "queue", "ame", "Name of the queue to be used!")
	flag.UintVar(&IDshift, "idshift", 0, "Value of the starting ID of the workers")
	flag.StringVar(&configuration, "c", "", "Configuration file for Ame")
	flag.StringVar(&DownloadDirectory, "download", "/tmp", "Directory for storage of temporary files.")
	flag.StringVar(&StorageDirectory, "storage", "/tmp", "Directory for storage of download files.")
	flag.Float64Var(&PersistenceTime, "persistence", 30, "Time zip files will stick around in minutes")
	flag.Float64Var(&CleanTime, "cleanperiod", 10, "Period for files to be deleted")
	flag.Parse()

	ctx = context.Background()
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

	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalln(err)
	}

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
	// Initialize Redis

	redisClient := redis.NewClient(&redis.Options{
		Addr:     REDISURI + ":" + strconv.Itoa(int(REDISIP)),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

	// Start Cleaners

	cleaner := Cleaner{}
	cleaner.redisClient = redisClient
	cleaner.trashCan = new([]Trash)

	go cleaner.Start()
	// Start workers

	for id := uint(1); id <= workercount; id++ {

		newworker := Worker{id + IDshift,
			channel,
			queue.Name,
			redisClient,
			ameinstance,
			&cleaner}

		go newworker.Work()
	}

	forever := make(chan bool)
	<-forever
}
