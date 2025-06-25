package api

import (
	"net/http"
	"time"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
)

func handleTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		task.Date = nextDate
		if err := db.UpdateTask(task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
