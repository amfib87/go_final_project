package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func signinHandler(res http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer
	var answer answerSignin
	var pass typePass

	passVar := os.Getenv("TODO_PASSWORD")
	if len(passVar) == 0 {
		writeJson(res, answer)
		return
	}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &pass); err != nil {
		answer.Error = err.Error()
		writeJson(res, answer)
		return
	}

	if pass.Password != passVar {
		answer.Error = fmt.Sprint("wrong password")
		writeJson(res, answer)
		return
	}

	hashPass := sha256.Sum256([]byte(passVar))
	str := fmt.Sprintf("%x", hashPass)

	claims := jwt.MapClaims{
		"hash": str,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(secret)) // получаем подписанный токен
	if err != nil {
		answer.Error = err.Error()
		writeJson(res, answer)
	}

	answer.Token = signedToken
	writeJson(res, answer)

}
