package main

import (
	"log/slog"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/teryble09/17.07.2025/internal/archiver/handler"
	"github.com/teryble09/17.07.2025/internal/archiver/service"
	"github.com/teryble09/17.07.2025/internal/config"
	"github.com/teryble09/17.07.2025/internal/storage"
	"golang.org/x/sync/semaphore"
)

func main() {
	_, curPath, _, _ := runtime.Caller(0)
	curPath, found := strings.CutSuffix(curPath, "cmd/archiver/main.go")
	if found != true {
		panic("Can not cut suffix to get to config: " + curPath)
	}
	cfg := config.MustLoad(curPath + "config.yaml")

	logger := slog.Default()
	sem := semaphore.NewWeighted(int64(cfg.MaxCurrentTasks))
	stor := storage.NewInMemoryStorage(cfg.MaxURLsInTask)

	srv := service.TaskService{
		Cfg:       cfg,
		Logger:    logger,
		Semaphore: sem,
		Storage:   stor,
	}

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", handler.CreateTask(srv))
		r.Post("/{task_id}/urls", handler.AddUrlToTask(srv))
		r.Get("/{task_id}/", handler.GetStatus(srv))
	})

	logger.Info("Starting http server on port " + cfg.Port)
	err := http.ListenAndServe("0.0.0.0:"+cfg.Port, r)
	if err != nil {
		panic(err.Error())
	}
}
