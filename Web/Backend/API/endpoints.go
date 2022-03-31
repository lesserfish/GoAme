package main

import "net/http"

func (server Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
