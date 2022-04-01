package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/streadway/amqp"
)

type Server struct {
	Router         *mux.Router
	DB             *sql.DB
	AMQPConnection *amqp.Connection
	AMQPChannel    *amqp.Channel
	RedisClient    *redis.Client
	queueName      string
}
type InitOptions struct {
	DB        string
	amqpADDR  string
	amqpPORT  string
	redisADDR string
	redisPORT string
	redisProc string
	queue     string
}

func CreateServer(options InitOptions) (*Server, error) {

	db, err := sql.Open("sqlite3", options.DB)
	if err != nil {
		return nil, err
	}
	amqpaddr := options.amqpADDR + ":" + options.amqpPORT
	amqpconnection, err := amqp.Dial(amqpaddr)

	if err != nil {
		return nil, err
	}

	amqpchannel, err := amqpconnection.Channel()
	if err != nil {
		log.Println(err)
	}

	amqpchannel.QueueDeclare(
		options.queue,
		false,
		false,
		false,
		false,
		nil)

	redisclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	server := Server{mux.NewRouter(),
		db,
		amqpconnection, amqpchannel,
		redisclient,
		options.queue}
	return &server, nil
}
func (server Server) Initiate() {
	server.CreateHandlers()

	_, err := server.DB.Exec("CREATE TABLE IF NOT EXISTS clients (ip TEXT NOT NULL, date TEXT, reqsize INTEGER);")

	if err != nil {
		log.Fatalln("Failed to create table. Error: " + err.Error())
	}
}
func (server Server) Serve(addr string) {
	http.Handle("/", server.Router)
	log.Println("Starting server at " + addr)

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Println(err)
	}
}
func (server Server) CreateHandlers() {
	server.Router.HandleFunc("/post", Wrap(server.PostHandler, server.Logger, server.Authorize, server.CheckPostValidity, server.RegisterRequest))
	server.Router.HandleFunc("/get", Wrap(server.GetHandler, server.Logger))
	server.Router.HandleFunc("/help", Wrap(server.HelpHandler, server.Logger))
}
func (server Server) Close() {
	server.DB.Close()
	server.AMQPConnection.Close()
}
