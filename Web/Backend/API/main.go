package main

import (
	"flag"
	"log"
	"strconv"

	ame "github.com/lesserfish/GoAme/Ame"
)

var (
	permanence  uint64
	maxrequests uint64
)

type PostStruct struct {
	AmeInput ame.Input
}

func main() {
	var address string
	var port uint64
	var database string

	flag.StringVar(&database, "db", "/tmp/db.sqlite", "path of sqlite3 database")
	flag.StringVar(&address, "addr", "localhost", "ip address of host")
	flag.Uint64Var(&port, "p", 9000, "port of where to serve")
	flag.Uint64Var(&permanence, "permanence", 240, "Duration each request will be stored in memory")
	flag.Uint64Var(&maxrequests, "maxreq", 5000, "Maximum number of requests per client")
	flag.Parse()

	addr := address + ":" + strconv.Itoa(int(port))
	options := InitOptions{
		DB: database}

	server, err := CreateServer(options)

	if err != nil {
		log.Fatalln("Failed to create server. Error: " + err.Error())
	}

	server.Initiate()
	server.Serve(addr)
	server.Close()
}
