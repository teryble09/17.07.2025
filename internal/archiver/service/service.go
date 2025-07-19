package service

import (
	"log/slog"

	"github.com/teryble09/17.07.2025/internal/archiver/model"
	"github.com/teryble09/17.07.2025/internal/config"
	"golang.org/x/sync/semaphore"
)

type TaskService struct {
	Cfg       *config.Config
	Logger    *slog.Logger
	Storage   model.TaskRepository
	Semaphore *semaphore.Weighted
}
