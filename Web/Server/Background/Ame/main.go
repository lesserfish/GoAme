package main

import (
    "context"
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "strconv"
    "time"

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
    ToolsDirectory string
    StorageDirectory  string
    PersistenceTime   float64
    CleanTime         float64
    ctx               context.Context
)

var amqpConnection *amqp.Connection;
var amqpChannel *amqp.Channel;

func generateRabbitInstance() error {
    log.Println("Generating Rabbit instance")
    fulladdr := AMQPURI + ":" + strconv.Itoa(int(AMQPIP))

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
        queuename,
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
        Addr:     REDISURI + ":" + strconv.Itoa(int(REDISIP)),
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

var ameInstance *ame.AmeKanji
func generateAmeInstance() error {
    log.Println("Generating Ame instance")
    config_content, err := ioutil.ReadFile(configuration)

    if err != nil {
         return err
    }

    var config ame.Configuration
    json.Unmarshal(config_content, &config)

    ameInstance, err = ame.Initialize(config)

    if err != nil {
        return err
    }
    return nil
}
func getAmeInstance() *ame.AmeKanji {
    if(ameInstance == nil){
        err := generateAmeInstance()
        if(err != nil){
            log.Fatalln("ERROR: Failed to initialize Ame. CRITICAL FAILURE!")
        }
        return ameInstance
    }

    return ameInstance
}
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
    flag.StringVar(&ToolsDirectory, "tools", "/tmp", "Directory for storage of resources.")
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


    generateRabbitInstance()
    generateAmeInstance()
    generateRedisInstance()

    // Start Cleaners

    cleaner := Cleaner{}
    cleaner.trashCan = new([]Trash)
    go cleaner.Start()

    // Start workers

    for id := uint(1); id <= workercount; id++ {

        newworker := Worker{id + IDshift,
        queuename,
        &cleaner}

        go newworker.Work()
    }

    forever := make(chan bool)
    <-forever
}
