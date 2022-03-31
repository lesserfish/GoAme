package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Files struct {
	UUID          uuid.UUID
	creation_time time.Time
}

var SavedFiles []Files

type Cleaner struct {
	redisClient redis.Conn
}

func (cleaner Cleaner) DeleteFileAt(position int) {
	SavedFiles[position] = SavedFiles[len(SavedFiles)-1]
	SavedFiles = SavedFiles[:len(SavedFiles)-1]
}

func (cleaner Cleaner) Clean() {
	for id, File := range SavedFiles {
		now := time.Now()
		diff := now.Sub(File.creation_time)
		if diff.Minutes() >= PersistenceTime {
			path := DownloadDirectory + "/" + GetZipnameFromID(File.UUID)
			cleaner.ReportDeleted(File.UUID)
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
		case <-time.After(time.Duration(CleanTime) * time.Minute):
			cleaner.Clean()
		}
	}
}

func (cleaner Cleaner) AddFile(id uuid.UUID, created_at time.Time) {
	SavedFiles = append(SavedFiles, Files{id, created_at})
}
