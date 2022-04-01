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

func (server Server) HelpHandler(rw http.ResponseWriter, r *http.Request) {
	kanji := r.FormValue("kanji")
	if kanji == "" {
		ErrorResponse(rw, "Failed to specify Kanji", http.StatusBadRequest)
	}

	smtm, err := server.DB.Prepare("SELECT kana FROM kanjikana where kanji == ?;")

	if err != nil {
		ErrorResponse(rw, "Internal error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	KKRow := []string{}
	rows, err := smtm.Query(kanji)

	if err != nil {
		ErrorResponse(rw, "Internal Error", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var kanjireading string
		rows.Scan(&kanjireading)

		KKRow = append(KKRow, kanjireading)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusUnsupportedMediaType)

	response := struct {
		Message  string
		Response struct {
			Kanji string
			Kana  []string
		}
	}{}

	response.Message = "OK"
	response.Response = struct {
		Kanji string
		Kana  []string
	}{kanji, KKRow}

	byteresponse, _ := json.Marshal(response)
	rw.Write(byteresponse)

}
func GetZipnameFromID(id uuid.UUID) string {
	out := "out_" + id.String() + ".zip"
	return out
}
