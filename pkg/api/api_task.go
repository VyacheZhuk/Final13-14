package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handlePutTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		writeJSON(w, map[string]string{"error": "Ошибка"}, http.StatusMethodNotAllowed)

	}
}
func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
	const dateLayout = "20060102"

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Неправильный формат данных"}, http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	if len(task.Date) != len(dateLayout) { // Используем длину константы для проверки
		writeJSON(w, map[string]string{"error": "Некорректный формат даты"}, http.StatusBadRequest)
		return
	}

	_, err := time.Parse(dateLayout, task.Date) // Используем константу для парсинга
	if err != nil {
		writeJSON(w, map[string]string{"error": "Некорректная дата"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Заголовок не может быть пустым"}, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		parts := strings.Fields(task.Repeat)
		if len(parts) != 2 {
			writeJSON(w, map[string]string{"error": "Некорректное правило повтора"}, http.StatusBadRequest)
			return
		}
		_, err := strconv.Atoi(parts[1])
		if err != nil {
			writeJSON(w, map[string]string{"error": "Некорректное правило повтора"}, http.StatusBadRequest)
			return
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Задача успешно обновлена"}, http.StatusOK)
}
