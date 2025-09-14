package api

import (
	"bytes"
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type Answer struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var task db.Task
	var answer Answer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		answer.Error = err.Error()
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		answer.Error = err.Error()
	}

	if task.Title == "" {
		answer.Error = "title is empty"
	}

	err = checkDate(&task)
	if err != nil {
		answer.Error = err.Error()
	}

	answer.Id = task.ID
	writeJson(res, answer, http.StatusOK)

}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(formatDate)
	}

	t, err := time.Parse(formatDate, task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	// если сегодня (now) больше task.Date (t)
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(formatDate)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		return err
	}
	task.ID = strconv.Itoa(int(id))

	return nil

}

func writeJson(res http.ResponseWriter, data any, code int) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(code)
	res.Write([]byte(resp))

}
