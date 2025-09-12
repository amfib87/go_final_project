package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func doneHandler(res http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")

	var answer Answer
	task, err := db.GetTask(id)
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer)
			return
		}
	} else {
		now := time.Now()

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer)
			return
		}

		err = db.UpdateDate(next, task.ID)
		if err != nil {
			answer.Error = err.Error()
			writeJson(res, answer)
			return
		}
	}

	writeJson(res, answer)

}
