package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const limit int = 50

func tasksHandler(res http.ResponseWriter, req *http.Request) {
	condText := req.URL.Query().Get("search")

	var condDate string
	if condText != "" {
		date, err := time.Parse("02.01.2006", condText)
		if err == nil {
			condDate = date.Format(formatDate)
		}
	}

	var answer Answer
	tasks, err := db.Tasks(limit, condText, condDate) // в параметре максимальное количество записей
	if err != nil {

		answer.Error = err.Error()

		writeJson(res, answer, http.StatusBadRequest)
		return
	}
	writeJson(res, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}
