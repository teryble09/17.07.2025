package handler

import (
	"encoding/json"
	"net/http"

	"github.com/teryble09/17.07.2025/internal/archiver"
)

func CreateTask(app archiver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := app.Semaphore.TryAcquire(1)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		ID := app.Storage.CreateTask()
		jsonID, err := json.Marshal(ID)
		if err != nil {
			app.Logger.Error("Could not marshall into json", "id", ID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonID)
	}
}
