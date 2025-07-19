package repository

import (
	"errors"

	"github.com/teryble09/17.07.2025/internal/archiver/model"
)

type TaskRepository interface {
	CreateTask() model.TaskID
	AddURL(model.TaskID, string) error
	Status(model.TaskID) ([]model.Url, error)
	ChangeStatus(id model.TaskID, url string, newStatus string) error
	LoadArchive(model.TaskID) ([]byte, error)
	WriteToArchive(id model.TaskID, filename []byte, file []byte) (archiveFinished bool, err error)
}

var (
	ErrTaskNotFound             = errors.New("task not found")
	ErrMaximumTaskNumberReached = errors.New("maximum task number reached")
	ErrArchiveNotReady          = errors.New("archive not ready")
	ErrFailedWrite              = errors.New("error writing")
)
