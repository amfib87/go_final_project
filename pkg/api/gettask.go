package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

func getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var answer Answer
	task, err := db.GetTask(id)
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer, http.StatusBadRequest)
		return
	}

	writeJson(res, task, http.StatusOK)
}
