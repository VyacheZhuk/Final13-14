package api

import (
	"net/http"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, TasksResponse{Tasks: tasks}, http.StatusOK)
}
