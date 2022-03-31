package main

import (
	"time"

	"github.com/google/uuid"
)

type Files struct {
	UUID          uuid.UUID
	creation_time time.Time
}

var SavedFiles []Files

func DeleteFileAt(position int) {
	SavedFiles[position] = SavedFiles[len(SavedFiles)-1]
	SavedFiles = SavedFiles[:len(SavedFiles)-1]
}

func Clean() {
	for id, File := range SavedFiles {
		now := time.Now()
		diff := now.Sub(File.creation_time)

		if diff.Minutes() >= PersistenceTime {
			path := DownloadDirectory + "/out_" + File.UUID.String() + ".zip"
			RemoveFile(path)
			DeleteFileAt(id)
			Clean()
			return
		}
	}
}

func CleanTasker() {

	for {
		select {
		case <-time.After(time.Duration(CleanTime) * time.Minute):
			Clean()
		}
	}
}

func AddFile(id uuid.UUID, created_at time.Time) {
	SavedFiles = append(SavedFiles, Files{id, created_at})
}
