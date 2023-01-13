package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
)

type Server struct {
	Router         *mux.Router
	APIRouter      *mux.Router
	DB             *sql.DB
	publicdir      string
	bmondai        *bluemonday.Policy
}
type InitOptions struct {
	DB        string
	amqpADDR  string
	amqpPORT  string
	redisADDR string
	redisPORT string
	redisProc string
	queue     string
	publicdir string
}

func CreateServer(options InitOptions) (*Server, error) {

	db, err := sql.Open("sqlite3", options.DB)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	apirouter := router.PathPrefix("/api/").Subrouter()

    generateRedisInstance()
    generateRabbitInstance()

	bmonday := bluemonday.UGCPolicy()
	server := Server{router,
		apirouter,
		db,
		options.publicdir,
		bmonday,
	}
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
	server.APIRouter.HandleFunc("/post", Wrap(server.PostHandler, server.Logger, server.CORS, server.Authorize, server.CheckPostValidity, server.RegisterRequest)).Methods("POST", "OPTIONS")
	server.APIRouter.HandleFunc("/get", Wrap(server.GetHandler, server.Logger, server.CORS)).Methods("GET", "OPTIONS")
	server.APIRouter.HandleFunc("/help", Wrap(server.HelpHandler, server.Logger, server.CORS)).Methods("GET")
	server.APIRouter.HandleFunc("/help", Wrap(server.HelpHandler2, server.Logger, server.CORS)).Methods("POST", "OPTIONS")

	server.Router.PathPrefix("/static/").HandlerFunc(Wrap(http.StripPrefix("/static/", http.FileServer(http.Dir(server.publicdir))).ServeHTTP, server.CORS, server.Logger))
	server.Router.PathPrefix("/").HandlerFunc(Wrap(http.FileServer(http.Dir(server.publicdir)).ServeHTTP, server.CORS, server.Logger))
}
func (server Server) Close() {
	server.DB.Close()
}
