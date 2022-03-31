package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}
type InitOptions struct {
	DB string
}

func CreateServer(options InitOptions) (*Server, error) {

	db, err := sql.Open("sqlite3", options.DB)

	if err != nil {
		return nil, err
	}

	server := Server{mux.NewRouter(), db}
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
	server.Router.HandleFunc("/", Wrap(server.PostHandler, server.Logger, server.Authorize, server.CheckPostValidity, server.RegisterRequest))
}
func (server Server) Close() {
	server.DB.Close()
}
