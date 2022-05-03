package main

import (
	"context"
	"flag"
	"log"
	"strconv"
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

func main() {
	var address string
	var port uint64
	var database string
	var amqpaddr string
	var amqpport uint64
	var redisaddr string
	var redisport uint64
	var redisproc string
	var queue string
	var publicdir string

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
	flag.StringVar(&corsoriginpolicy, "corsorigin", "", "Cors policy")
	flag.StringVar(&corsmethodpolicy, "corsmethod", "*", "Cors policy")
	flag.StringVar(&corsheaderpolicy, "corsheader", "*", "Cors policy")
	flag.StringVar(&publicdir, "publicdir", "/tmp", "Directory of public HTML files.")
	flag.Parse()

	ctx = context.Background()

	addr := address + ":" + strconv.Itoa(int(port))
	options := InitOptions{
		DB:        database,
		amqpADDR:  amqpaddr,
		amqpPORT:  strconv.Itoa(int(amqpport)),
		redisADDR: redisaddr,
		redisPORT: strconv.Itoa(int(redisport)),
		redisProc: redisproc,
		queue:     queue,
		publicdir: publicdir}

	server, err := CreateServer(options)

	if err != nil {
		log.Fatalln("Failed to create server. Error: " + err.Error())
	}

	server.Initiate()
	server.Serve(addr)
	server.Close()
}
