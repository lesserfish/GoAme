package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	ame "github.com/lesserfish/GoAme/Ame"
	"github.com/streadway/amqp"
)

type PostStruct struct {
	AmeInput ame.Input
}
type Message struct {
	UUID  uuid.UUID
	Input ame.Input
}

func (server Server) PostHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ErrorResponse(rw, "Failed to parse body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	var postStruct PostStruct

	err = json.Unmarshal(body, &postStruct)

	if err != nil {
		ErrorResponse(rw, "Failed to parse body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	newid := uuid.New()
	message := Message{newid, postStruct.AmeInput}

	byteinfo, err := json.Marshal(message)

	if err != nil {
		ErrorResponse(rw, "Internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = server.AMQPChannel.Publish("", server.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteinfo})

	if err != nil {
		ErrorResponse(rw, "Internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusUnsupportedMediaType)

	response := make(map[string]string)
	response["message"] = "OK!"
	response["uuid"] = newid.String()

	byteresponse, _ := json.Marshal(response)
	rw.Write(byteresponse)
}

func (server Server) GetHandler(rw http.ResponseWriter, r *http.Request) {
	reqid := r.FormValue("id")
	if reqid == "" {
		ErrorResponse(rw, "Failed to specify id", http.StatusBadRequest)
	}

	redisout := server.RedisClient.HGetAll(ctx, reqid)

	result, err := redisout.Result()

	if err != nil {
		ErrorResponse(rw, "Internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if len(result) == 0 {
		ErrorResponse(rw, "ID not found", http.StatusBadRequest)
		log.Println(err)
		return
	}
	status := result["Status"]
	progress := result["Progress"]

	if status == "Success" {
		filename := DownloadDirectory + "/" + GetZipnameFromID(uuid.MustParse(reqid))
		rw.Header().Add("content-disposition", "filename=\"out.zip\"")
		http.ServeFile(rw, r, filename)
	} else {

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnsupportedMediaType)

		response := make(map[string]string)
		response["Status"] = status
		response["Progress"] = progress

		byteresponse, _ := json.Marshal(response)
		rw.Write(byteresponse)
	}

}

func GetZipnameFromID(id uuid.UUID) string {
	out := "out_" + id.String() + ".zip"
	return out
}
