package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Wrap(endpoint http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for id := len(middleware) - 1; id >= 0; id-- {
		endpoint = middleware[id](endpoint)
	}
	return endpoint
}

func (server Server) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handlers.LoggingHandler(log.Writer(), next).ServeHTTP(rw, r)
	}
}
func (server Server) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		addCorsHeader(rw)

		if r.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		} else {
			handlers.CORS()(next).ServeHTTP(rw, r)
		}
	}
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
}

func (server Server) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		remote_addr := ""
		forwarded := r.Header.Get("X-FORWARDED-FOR")
		if forwarded != "" {
			remote_addr = forwarded
		} else {
			remote_addr = r.RemoteAddr
		}

		smtm, err := server.DB.Prepare("SELECT date, reqsize FROM clients where IP == ?;")

		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}

		rows, err := smtm.Query(remote_addr)

		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}

		reqsum := uint64(0)
		for rows.Next() {
			var datestr string
			var reqsize int
			rows.Scan(&datestr, &reqsize)

			date, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", datestr)

			if err != nil {
				ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
				return
			}

			now := time.Now()

			diff := now.Sub(date)

			if diff.Minutes() <= float64(permanence) {
				reqsum += uint64(reqsize)
			}
		}

		if reqsum >= maxrequests {
			log.Println("Remote address " + remote_addr + " has exceeded it's quota.")
			ErrorResponse(rw, "Too many requests!", http.StatusTooManyRequests)
			return
		} else {
			log.Println("Remote address " + remote_addr + " has a total of " + strconv.Itoa(int(reqsum)) + " requests.")
		}

		next(rw, r)
	}
}

func (server Server) CheckPostValidity(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			ErrorResponse(rw, "Content Type is not JSON", http.StatusUnsupportedMediaType)
			return
		}

		var request PostStruct
		body, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(body, &request)

		if err != nil {
			ErrorResponse(rw, "Failed to parse body", http.StatusBadRequest)
			return
		}

		if len(request.AmeInput.Input) == 0 {
			ErrorResponse(rw, "Empty Request", http.StatusBadRequest)
			return
		}

		if len(request.AmeInput.Input) > int(maxrequests) {
			ErrorResponse(rw, "Request too large", http.StatusBadRequest)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		next(rw, r)
	}
}
func ErrorResponse(rw http.ResponseWriter, message string, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	response := make(map[string]string)
	response["Message"] = message

	byteresponse, _ := json.Marshal(response)

	rw.Write(byteresponse)
}
func (server Server) RegisterRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		remote_addr := ""
		forwarded := r.Header.Get("X-FORWARDED-FOR")
		if forwarded != "" {
			remote_addr = forwarded
		} else {
			remote_addr = r.RemoteAddr
		}

		time := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")

		var request PostStruct

		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)

		if err != nil {
			ErrorResponse(rw, "Failed to parse body", http.StatusBadRequest)
			log.Println(err)
			return
		}

		reqsize := len(request.AmeInput.Input)

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		tx, err := server.DB.Begin()
		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}
		smtm, err := tx.Prepare("INSERT INTO clients VALUES (?, ?, ?)")

		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}

		_, err = smtm.Exec(remote_addr, time, reqsize)

		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()

		if err != nil {
			ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
			return
		}

		log.Println("Registerd a request of size " + strconv.Itoa(reqsize) + " from address " + remote_addr)
		next(rw, r)
	}
}
