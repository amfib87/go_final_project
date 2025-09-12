package main

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	fileDb := "scheduler.db"
	fileDbEnv := os.Getenv("TODO_DBFILE")

	if fileDbEnv != "" {
		fileDb = fileDbEnv
	}

	err_db := db.Init(fileDb)
	if err_db != nil {
		panic(err_db)
	}
	defer db.Db.Close()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()

	port := ":" + os.Getenv("TODO_PORT")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
