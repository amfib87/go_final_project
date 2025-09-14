package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func doneHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "incorrect method od request", http.StatusBadRequest)
	}

	id := req.URL.Query().Get("id")

	var answer Answer
	task, err := db.GetTask(id)
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer, http.StatusBadRequest)
			return
		}

		err = db.UpdateDate(next, task.ID)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer, http.StatusInternalServerError)
			return
		}
	}

	writeJson(res, answer, http.StatusOK)

}
