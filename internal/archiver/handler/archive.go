package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/teryble09/17.07.2025/internal/archiver/dto"
	"github.com/teryble09/17.07.2025/internal/archiver/service"
)

func GetArchive(srv service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetArchiveRequest
		id := chi.URLParam(r, "task_id")
		req.TaskId = id
		resp, err := srv.GetArchive(req)
		switch err {
		case service.ErrArchiveNotReady:
			srv.Logger.Warn("Trying to acces not completed archive")
			w.WriteHeader(http.StatusBadRequest)
			return
		case service.ErrTaskNotFound:
			srv.Logger.Warn("Trying to access to not existing task")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
		w.Header().Set("Content-Length", strconv.Itoa(len(resp.Data)))
		_, err = w.Write(resp.Data)
		if err != nil {
			srv.Logger.Error("Failed to send archive", "error", err)
		}
	}
}
