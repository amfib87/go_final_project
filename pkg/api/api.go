package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

const formatDate = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/signin", signinHandler)
}

var Port string

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(PassVar) > 0 {

			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			hashPass := sha256.Sum256([]byte(PassVar))
			// здесь код для валидации и проверки JWT-токена

			claims := jwt.MapClaims{}

			// парсим токен
			token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			hashTokenRaw := claims["hash"]
			hashTokenString := hashTokenRaw.(string)
			hashPassString := fmt.Sprintf("%x", hashPass)

			if hashPassString != hashTokenString {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if hashPass == [32]byte{0} {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

		}
		next(w, r)
	})
}
