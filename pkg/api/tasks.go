package api

import (
	"net/http"
	"strconv"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
)

type TaskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
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

	var responseTasks []TaskResponse
	for _, task := range tasks {
		responseTasks = append(responseTasks, TaskResponse{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	if responseTasks == nil {
		responseTasks = []TaskResponse{}
	}

	writeJSON(w, TasksResponse{Tasks: responseTasks}, http.StatusOK)
}
