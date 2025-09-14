package api

import (
	"bytes"
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func updateTaskHandler(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var task db.Task
	var answer Answer
	var status int

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		answer.Error = err.Error()
		status = http.StatusBadRequest
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		answer.Error = err.Error()
		status = http.StatusBadRequest
	}

	if task.Title == "" {
		answer.Error = "title is empty"
		status = http.StatusBadRequest
	}

	_, err = time.Parse(formatDate, task.Date)
	if err != nil {
		answer.Error = err.Error()
		status = http.StatusInternalServerError
	}

	now := time.Now()
	_, err = NextDate(now, task.Date, task.Repeat)
	if err != nil {
		answer.Error = err.Error()
		status = http.StatusBadRequest
	}

	err = db.UpdateTask(&task)
	if err != nil {
		answer.Error = err.Error()
		status = http.StatusInternalServerError
	}

	if answer.Error == "" {
		status = http.StatusOK
	}
	writeJson(res, answer, status)

}
