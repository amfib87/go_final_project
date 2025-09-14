package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

const secret = "secretWord"

type answerSignin struct {
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

type typePass struct {
	Password string `json:"password"`
}

var PassVar string

func signinHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "incorrect method od request", http.StatusBadRequest)
	}

	var buf bytes.Buffer
	var answer answerSignin
	var pass typePass

	if len(PassVar) == 0 {
		writeJson(res, answer, http.StatusInternalServerError)
		return
	}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &pass); err != nil {
		answer.Error = err.Error()
		writeJson(res, answer, http.StatusBadRequest)
		return
	}

	if pass.Password != PassVar {
		answer.Error = fmt.Sprint("wrong password")
		writeJson(res, answer, http.StatusBadRequest)
		return
	}

	hashPass := sha256.Sum256([]byte(PassVar))
	str := fmt.Sprintf("%x", hashPass)

	claims := jwt.MapClaims{
		"hash": str,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(secret)) // получаем подписанный токен
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer, http.StatusInternalServerError)
	}

	answer.Token = signedToken
	writeJson(res, answer, http.StatusOK)

}
