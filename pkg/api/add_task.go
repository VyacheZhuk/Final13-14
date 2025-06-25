package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
)

const dateLayout = "20060102"

type taskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateAndAdjustTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, taskResponse{ID: id}, http.StatusCreated)
}

func validateAndAdjustTask(task *db.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("task title cannot be empty")
	}

	now := time.Now()
	today := now.Format(dateLayout)

	if task.Date == "" || task.Date == "today" {
		task.Date = today
		return nil
	}

	parsedDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	if isBeforeToday(now, parsedDate) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = nextDate
		}
	}

	return nil
}

func isBeforeToday(now, date time.Time) bool {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return checkDate.Before(today)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	writeJSON(w, taskResponse{Error: errorMsg}, statusCode)
}
