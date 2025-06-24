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
	var input struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, map[string]string{"error": "Неправильный формат данных"}, http.StatusBadRequest)
		return
	}

	// Проверка ID
	if input.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Проверка даты (формат YYYYMMDD)
	if len(input.Date) != 8 {
		writeJSON(w, map[string]string{"error": "Некорректный формат даты"}, http.StatusBadRequest)
		return
	}

	// Проверка что дата валидна
	_, err := time.Parse("20060102", input.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Некорректная дата"}, http.StatusBadRequest)
		return
	}

	// Проверка заголовка
	if input.Title == "" {
		writeJSON(w, map[string]string{"error": "Заголовок не может быть пустым"}, http.StatusBadRequest)
		return
	}

	// Проверка правила повтора (если указано)
	if input.Repeat != "" {
		parts := strings.Fields(input.Repeat)
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

	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат идентификатора"}, http.StatusBadRequest)
		return
	}

	task := db.Task{
		ID:      id,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Задача успешно обновлена"}, http.StatusOK)
}
