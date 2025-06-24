package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DataFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	date, err := time.Parse(DataFormat, dstart)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid days format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid days interval")
		}
		return calculateDailyRepeat(now, date, days), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly format")
		}
		return calculateYearlyRepeat(now, date), nil

	default:
		return "", errors.New("unsupported repeat format")
	}
}

func calculateDailyRepeat(now, date time.Time, days int) string {
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DataFormat)
}

func calculateYearlyRepeat(now, date time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DataFormat)
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if now == "" || date == "" || repeat == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	nowTime, err := time.Parse(DataFormat, now)
	if err != nil {
		http.Error(w, "Invalid now parameter", http.StatusBadRequest)
		return
	}

	result, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
