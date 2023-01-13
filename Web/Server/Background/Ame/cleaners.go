package main

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Trash struct {
	UUID          uuid.UUID
	creation_time time.Time
}

type Cleaner struct {
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
			path := StorageDirectory + "/" + GetZipnameFromID(File.UUID)
			cleaner.ReportDeleted(File.UUID)
			err := RemoveFile(path)
			if err != nil {
				log.Println("Failed to delete file. Error: " + err.Error())
			}
			log.Println("[Cleaner] Deleted file with id: " + File.UUID.String())
			cleaner.DeleteFileAt(id)
			cleaner.Clean()
			return
		}
	}
}
func (cleaner Cleaner) ReportDeleted(id uuid.UUID) {
	getRedisInstance().HMSet(ctx, id.String(),
		"Status", "Deleted",
		"Progress", "1")
}
func (cleaner Cleaner) Start() {

	log.Println("[Cleaner] Starting work!")
	for {
		select {
		case <-time.After(time.Duration(CleanTime) * time.Minute):
			cleaner.Clean()
		}
	}
}

func (cleaner Cleaner) AddTrash(id uuid.UUID, created_at time.Time) {
	log.Println("[Cleaner] Added file to trashcan with id: " + id.String())
	*cleaner.trashCan = append(*cleaner.trashCan, Trash{id, created_at})
}
