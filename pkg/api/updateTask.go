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

	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		answer.Error = err.Error()
	}

	now := time.Now()
	_, err = NextDate(now, task.Date, task.Repeat)
	if err != nil {
		answer.Error = err.Error()
	}

	err = db.UpdateTask(&task)
	if err != nil {
		answer.Error = err.Error()
	}

	writeJson(res, answer)

}
