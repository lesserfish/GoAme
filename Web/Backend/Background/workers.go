package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	ame "github.com/lesserfish/GoAme/Ame"
	"github.com/streadway/amqp"
)

type Message struct {
	UUID  uuid.UUID
	Input ame.Input
}

type Worker struct {
	workerID    uint
	channel     *amqp.Channel
	queueName   string
	redisClient redis.Conn
	AmeKanji    *ame.AmeKanji
}

type jstest struct {
	Name []string
	Age  int
}

func (worker Worker) Work() {
	logmsg := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] Starting work!"
	log.Println(logmsg)

	msgs, err := worker.channel.Consume(
		worker.queueName,
		strconv.Itoa(int(worker.workerID)),
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		errmsg := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] "
		errmsg += "Error: " + err.Error()
		log.Println(errmsg)
		return
	}

	for msg := range msgs {
		var message Message
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Println("Error parsing message. Error: " + err.Error())
			msg.Ack(false)
			continue
		}
		worker.AcceptRequest(message.UUID)

		// Create directory for request

		new_directory := StorageDirectory + "/" + "dir_" + message.UUID.String() + "/"
		new_media_directory := new_directory + "Media" + "/"

		_, err = os.Stat(new_directory)

		if !os.IsNotExist(err) {
			msg.Ack(false)
			worker.ReportError(message.UUID)
			log.Println("Directory " + new_directory + " already existed. Critical failure!")
			continue
		}

		err1 := CreateDir(new_directory)
		err2 := CreateDir(new_media_directory)

		if err1 != nil || err2 != nil {
			msg.Ack(false)
			worker.ReportError(message.UUID)
			log.Println("Failed creating directory " + new_directory + ". Error: " + err.Error())
			continue
		}

		for i, _ := range message.Input.Input {
			message.Input.Input[i]["savepath"] = new_media_directory
		}

		deckfile := new_directory + "anki_deck.txt"

		// Invoke AmeKanji
		worker.AmeKanji.URenderAndSave(message.Input, deckfile, func(p float64) {
			worker.LogProgress(message.UUID, p)
		})

		// Create zip file and move it to Download directory

		zipdir := DownloadDirectory + "/" + "out_" + message.UUID.String() + ".zip"
		err = ZipDir(new_directory, zipdir)

		if err != nil {
			msg.Ack(false)
			worker.ReportError(message.UUID)
			log.Println("Could not create ZIP file. Error: " + err.Error())
			continue
		}

		// Delete previously create directory

		RemoveDir(new_directory)

		// Success
		worker.ReportSuccess(message.UUID)
		msg.Ack(false)
	}

	text := "[Worker " + strconv.Itoa(int(worker.workerID)) + "] " + "Exiting!"
	fmt.Println(text)

}

func (worker Worker) AcceptRequest(id uuid.UUID) {
	worker.redisClient.Do("HMSET", id.String(),
		"Status", "Accepted",
		"Progress", "0")
}
func (worker Worker) LogProgress(id uuid.UUID, progress float64) {
	worker.redisClient.Do("HMSET", id.String(),
		"Status", "In Progress",
		"Progress", fmt.Sprint(progress))
}
func (worker Worker) ReportError(id uuid.UUID) {
	worker.redisClient.Do("HMSET", id.String(),
		"Status", "Failed",
		"Progress", "0")
}
func (worker Worker) ReportSuccess(id uuid.UUID) {
	worker.redisClient.Do("HMSET", id.String(),
		"Status", "Success",
		"Progress", "1")
}
func CreateDir(path string) error {
	err := os.Mkdir(path, 0755)
	return err
}
func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	return err
}
func ZipDir(source string, target string) error {
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := zip.NewWriter(out)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
