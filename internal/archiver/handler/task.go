package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/teryble09/17.07.2025/internal/archiver/dto"
	"github.com/teryble09/17.07.2025/internal/archiver/service"
)

func CreateTask(srv service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := srv.Semaphore.TryAcquire(1)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		ID := srv.Storage.CreateTask()
		jsonID, err := json.Marshal(ID)
		if err != nil {
			srv.Logger.Error("Could not marshall into json", "id", ID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonID)
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

		var dto dto.AddURLRequest
		err = json.Unmarshal(data, &dto)
		if err != nil {
			srv.Logger.Warn("Could not unmarshall AddURLRequest", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
