package service

import (
	"log/slog"

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
