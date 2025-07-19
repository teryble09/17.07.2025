package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teryble09/17.07.2025/internal/archiver/dto"
	"github.com/teryble09/17.07.2025/internal/archiver/service"
)

func CreateTask(srv service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		resp, err := srv.CreateTask(dto.CreateTaskRequest{})
		if err == service.ErrServerBusy {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			srv.Logger.Error("Could not marshall into json", "id", resp.Id)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResp)
		if err != nil {
			srv.Logger.Error("Could not write jsonID into ResponseWriter", "error", err.Error())
		}
	}
}

func AddUrlToTask(srv service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			srv.Logger.Error("Could not read request body", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req dto.AddURLRequest
		id := chi.URLParam(r, "task_id")
		req.TaskId = id

		err = json.Unmarshal(data, &req)
		if err != nil {
			srv.Logger.Warn("Could not unmarshall AddURLRequest", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = srv.AddURL(req)
		if err != nil {
			switch err {
			case service.ErrMaximumTaskNumberReached:
				srv.Logger.Warn("Trying to exceed max urls in task")
				w.WriteHeader(http.StatusBadRequest)
				return
			case service.ErrTaskNotFound:
				srv.Logger.Warn("Trying to access to not existing task")
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				srv.Logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetStatus(srv service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req dto.GetStatusRequest
		id := chi.URLParam(r, "task_id")
		req.TaskId = id

	}
}
