package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := Db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func Tasks(limit int, condText string, condDate string) ([]*Task, error) {
	var rowSelect string

	textLow := strings.ToLower(condText)
	textUp := strings.ToUpper(condText)

	textLow = "%" + textLow + "%"
	textUp = "%" + textUp + "%"

	if condDate != "" {
		rowSelect = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit`
	} else if condText != "" {
		rowSelect = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :textUp 
					OR title LIKE :textLow OR comment LIKE :textUp OR comment LIKE :textLow ORDER BY date LIMIT :limit`
	} else {
		rowSelect = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"
	}

	rows, err := Db.Query(rowSelect, sql.Named("limit", limit), sql.Named("date", condDate),
		sql.Named("textUp", textUp), sql.Named("textLow", textLow))
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	var tasks = []*Task{}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}

		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}

	return tasks, nil

}

func GetTask(id string) (task *Task, err error) {
	task = &Task{}

	rowSelect := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	err = Db.QueryRow(rowSelect, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil

}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET 	date    = :date,
								   	title   = :title,
								   	comment = :comment,
								   	repeat  = :repeat WHERE id = :id`
	res, err := Db.Exec(query, sql.Named("date", &task.Date), sql.Named("title", &task.Title), sql.Named("comment", &task.Comment),
		sql.Named("repeat", &task.Repeat), sql.Named("id", &task.ID))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := Db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := Db.Exec(query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
