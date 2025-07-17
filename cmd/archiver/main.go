package main

import (
	"log/slog"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/teryble09/17.07.2025/internal/config"
)

func main() {
	_, curPath, _, _ := runtime.Caller(0)
	curPath, found := strings.CutSuffix(curPath, "cmd/archiver/main.go")
	if found != true {
		panic("Can not cut suffix to get to config: " + curPath)
	}
	cfg := config.MustLoad(curPath + "config.yaml")

	logger := slog.Default()

	r := chi.NewRouter()

	logger.Info("Starting http server on port " + cfg.Port)
	http.ListenAndServe(cfg.Port, r)
}
