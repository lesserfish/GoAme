package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Files struct {
	UUID          uuid.UUID
	creation_time time.Time
}

type Cleaner struct {
	SavedFiles  []Files
	redisClient redis.Conn
}

func (cleaner Cleaner) DeleteFileAt(position int) {
	fmt.Println("Deleting file")
	cleaner.SavedFiles[position] = cleaner.SavedFiles[len(cleaner.SavedFiles)-1]
	cleaner.SavedFiles = cleaner.SavedFiles[:len(cleaner.SavedFiles)-1]
}

func (cleaner Cleaner) Clean() {
	log.Println("Cleaning " + strconv.Itoa(len(cleaner.SavedFiles)) + " files")
	for id, File := range cleaner.SavedFiles {
		now := time.Now()
		diff := now.Sub(File.creation_time)

		fmt.Println(PersistenceTime - diff.Minutes())
		if diff.Minutes() >= PersistenceTime {
			path := DownloadDirectory + "/out_" + File.UUID.String() + ".zip"
			cleaner.ReportDeleted(File.UUID)
			log.Panicln("Deleting files at " + path)
			err := RemoveFile(path)
			if err != nil {
				log.Println("Failed to delete file. Error: " + err.Error())
			}
			cleaner.DeleteFileAt(id)
			cleaner.Clean()
			return
		}
	}
}
func (cleaner Cleaner) ReportDeleted(id uuid.UUID) {
	cleaner.redisClient.Do("HMSET", id.String(),
		"Status", "Deleted",
		"Progress", "1")
}
func (cleaner Cleaner) CleanTasker() {

	for {
		select {
		case <-time.After(time.Duration(CleanTime) * time.Second):
			cleaner.Clean()
		}
	}
}

func (cleaner Cleaner) AddFile(id uuid.UUID, created_at time.Time) {
	fmt.Println("Adding file!")
	cleaner.SavedFiles = append(cleaner.SavedFiles, Files{id, created_at})
}
