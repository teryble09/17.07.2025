package model

import "errors"

type TaskRepository interface {
	CreateTask() TaskID
	AddURL(TaskID, string) error
	Status(TaskID) ([]Url, error)
	LoadArchive(TaskID) ([]byte, error)
	WriteToArchive(id TaskID, fuliname []byte, file []byte) error
}

var (
	ErrTaskNotFound             = errors.New("task not found")
	ErrMaximumTaskNumberReached = errors.New("maximum task number reached")
	ErrArchiveNotReady          = errors.New("archive not ready")
	ErrFailedWrite              = errors.New("error writing")
)
