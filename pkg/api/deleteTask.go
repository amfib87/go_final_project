package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {

	var answer Answer

	id := req.URL.Query().Get("id")

	err := db.DeleteTask(id)
	if err != nil {
		answer.Error = err.Error()
	}

	writeJson(res, answer, http.StatusOK)
}
