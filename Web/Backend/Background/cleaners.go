package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Trash struct {
	UUID          uuid.UUID
	creation_time time.Time
}

type Cleaner struct {
	redisClient redis.Conn
	trashCan    *[]Trash
}

func (cleaner Cleaner) DeleteFileAt(position int) {
	(*cleaner.trashCan)[position] = (*cleaner.trashCan)[len(*cleaner.trashCan)-1]
	*cleaner.trashCan = (*cleaner.trashCan)[:len(*cleaner.trashCan)-1]
}

func (cleaner Cleaner) Clean() {
	for id, File := range *cleaner.trashCan {
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
func (cleaner Cleaner) Start() {
	for {
		select {
		case <-time.After(time.Duration(CleanTime) * time.Minute):
			cleaner.Clean()
		}
	}
}

func (cleaner Cleaner) AddTrash(id uuid.UUID, created_at time.Time) {
	*cleaner.trashCan = append(*cleaner.trashCan, Trash{id, created_at})
}
