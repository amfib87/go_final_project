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

	api.PassVar = os.Getenv("TODO_PASSWORD")
	if len(api.PassVar) == 0 {
		log.Print("password is empty")
	}

	api.Port = ":" + os.Getenv("TODO_PORT")
	if len(api.PassVar) == 0 {
		log.Print("port is empty")
	}

	db.FileDbEnv = os.Getenv("TODO_DBFILE")

}

func main() {
	fileDb := "scheduler.db"

	if db.FileDbEnv != "" {
		fileDb = db.FileDbEnv
	}

	err_db := db.Init(fileDb)
	if err_db != nil {
		panic(err_db)
	}
	defer db.Db.Close()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()

	err := http.ListenAndServe(api.Port, nil)
	if err != nil {
		panic(err)
	}
}
