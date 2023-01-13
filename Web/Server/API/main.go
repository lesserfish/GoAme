package main

import (
	"context"
	"flag"
	"log"
	"github.com/go-redis/redis/v8"
    "github.com/streadway/amqp"
    "strconv"
    "time"
)

var (
	permanence        uint64
	maxrequests       uint64
	maxhelprequests   uint64
	DownloadDirectory string
	ctx               context.Context
	corsoriginpolicy  string
	corsmethodpolicy  string
	corsheaderpolicy  string
)
var (
	address string
	port uint64
	database string
	amqpaddr string
	amqpport uint64
	redisaddr string
	redisport uint64
	redisproc string
	queue string
	publicdir string
)

var amqpConnection *amqp.Connection;
var amqpChannel *amqp.Channel;

func generateRabbitInstance() error {
    log.Println("Generating Rabbit instance")
    fulladdr := amqpaddr + ":" + strconv.Itoa(int(amqpport))

    var err error = nil
    amqpConnection, err = amqp.Dial(fulladdr)

    if(err != nil) {
        return err
    }

    amqpChannel, err = amqpConnection.Channel()
    if err != nil {
        return err
    }

    err = amqpChannel.Qos(1, 0, false)
    if err != nil {
        return err
    }
    _, err = amqpChannel.QueueDeclare(
        queue,
        false,
        false,
        false,
        false,
        nil,
    )

    if err != nil {
        return err
    }

    return nil
}
func getRabbitInstance() *amqp.Channel {

    if(amqpConnection == nil) {
        err := generateRabbitInstance()

        if(err != nil) {
            log.Println("ERROR: Could not create AMQP connection... Waiting 10 seconds and trying again.")
            time.Sleep(10 * time.Second)
            return getRabbitInstance()
        }

        return amqpChannel
    }

    if(amqpConnection.IsClosed()) {
        err := generateRabbitInstance()

        if(err != nil) {
            log.Println("ERROR: Could not create AMQP connection... Waiting 10 seconds and trying again.")
            time.Sleep(10 * time.Second)
            return getRabbitInstance()
        }

        return amqpChannel
    }

    return amqpChannel
}

var redisClient *redis.Client
func generateRedisInstance() error {
    log.Println("Generating Redis instance")
    redisClient = redis.NewClient(&redis.Options{
        Addr:     redisaddr + ":" + strconv.Itoa(int(redisport)),
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    return nil
}

func getRedisInstance() *redis.Client {
    if(redisClient == nil){
        generateRedisInstance()
    }
    return redisClient
}
func main() {


	flag.StringVar(&database, "db", "/tmp/db.sqlite", "path of sqlite3 database")
	flag.StringVar(&address, "addr", "localhost", "ip address of host")
	flag.Uint64Var(&port, "p", 9000, "port of where to serve")
	flag.Uint64Var(&permanence, "permanence", 240, "Duration each request will be stored in memory")
	flag.Uint64Var(&maxrequests, "maxreq", 5000, "Maximum number of requests per client")
	flag.Uint64Var(&maxhelprequests, "maxhelpreq", 1000, "Maximum number of help requests per client")
	flag.StringVar(&amqpaddr, "amqp", "amqp://localhost", "Address of amqp")
	flag.Uint64Var(&amqpport, "amqpport", 5672, "Amqp port")
	flag.StringVar(&queue, "queue", "ame", "Queue name")
	flag.StringVar(&redisaddr, "redis", "localhost", "Address of Redis")
	flag.Uint64Var(&redisport, "redisport", 6379, "Redis port")
	flag.StringVar(&redisproc, "redisproc", "tcp", "Proc of Redis")
	flag.StringVar(&DownloadDirectory, "download", "/tmp", "Directory for storage of download files.")
	flag.StringVar(&corsoriginpolicy, "corsorigin", "*", "Cors policy")
	flag.StringVar(&corsmethodpolicy, "corsmethod", "*", "Cors policy")
	flag.StringVar(&corsheaderpolicy, "corsheader", "*", "Cors policy")
	flag.StringVar(&publicdir, "publicdir", "/tmp", "Directory of public HTML files.")
	flag.Parse()

	ctx = context.Background()

	addr := address + ":" + strconv.Itoa(int(port))
	options := InitOptions{
		DB:        database,
		publicdir: publicdir}

	server, err := CreateServer(options)

	if err != nil {
		log.Fatalln("Failed to create server. Error: " + err.Error())
	}

	server.Initiate()
	server.Serve(addr)
	server.Close()
}
