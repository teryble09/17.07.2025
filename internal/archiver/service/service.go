package service

import (
	"errors"
	"log/slog"

	"github.com/teryble09/17.07.2025/internal/archiver/dto"
	"github.com/teryble09/17.07.2025/internal/archiver/model"
	"github.com/teryble09/17.07.2025/internal/archiver/repository"
	"github.com/teryble09/17.07.2025/internal/config"
	"golang.org/x/sync/semaphore"
)

type TaskService struct {
	Cfg       *config.Config
	Logger    *slog.Logger
	Storage   repository.TaskRepository
	Semaphore *semaphore.Weighted
}

var (
	ErrServerBusy               = errors.New("server is busy")
	ErrTaskNotFound             = errors.New("task not found")
	ErrUrlAlreadyExists         = errors.New("url already exists")
	ErrMaximumTaskNumberReached = errors.New("maximum task number reached")
	ErrArchiveNotReady          = errors.New("archive not ready")
)

func (srv *TaskService) CreateTask(req dto.CreateTaskRequest) (dto.CreateTaskResponse, error) {
	ok := srv.Semaphore.TryAcquire(1)
	if !ok {
		return dto.CreateTaskResponse{}, ErrServerBusy
	}

	TaskID := srv.Storage.CreateTask()
	return dto.CreateTaskResponse{Id: TaskID.Id}, nil
}

func (srv *TaskService) AddURL(req dto.AddURLRequest) (dto.AddURLResponse, error) {
	err := srv.Storage.AddURL(model.TaskID{Id: req.TaskId}, req.Adress)
	if err != nil {
		switch err {
		case repository.ErrUrlAlreadyExists:
			return dto.AddURLResponse{}, ErrUrlAlreadyExists
		case repository.ErrTaskNotFound:
			return dto.AddURLResponse{}, ErrTaskNotFound
		case repository.ErrMaximumTaskNumberReached:
			return dto.AddURLResponse{}, ErrMaximumTaskNumberReached
		default:
			return dto.AddURLResponse{}, errors.Join(errors.New("Could not add url to task"), err)
		}
	}

	go LoadFileAndArchive(srv, model.TaskID{Id: req.TaskId}, req.Adress)

	return dto.AddURLResponse{}, nil
}

func (srv *TaskService) GetStatus(req dto.GetStatusRequest) (dto.GetStatusResponse, error) {
	urls, err := srv.Storage.Status(model.TaskID{Id: req.TaskId})
	if err == repository.ErrTaskNotFound {
		return dto.GetStatusResponse{}, ErrTaskNotFound
	}
	return dto.GetStatusResponse{Urls: urls}, nil
}

func (srv *TaskService) GetArchive(req dto.GetArchiveRequest) (dto.GetArchiveResponse, error) {
	buf, err := srv.Storage.LoadArchive(model.TaskID{Id: req.TaskId})
	if err == repository.ErrArchiveNotReady {
		return dto.GetArchiveResponse{}, ErrArchiveNotReady
	}
	if err == repository.ErrTaskNotFound {
		return dto.GetArchiveResponse{}, ErrTaskNotFound
	}
	return dto.GetArchiveResponse{Data: buf}, nil
}
