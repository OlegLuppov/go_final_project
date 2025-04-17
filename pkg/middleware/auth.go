package middleware

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OlegLuppov/go_final_project/config"
	"github.com/golang-jwt/jwt/v5"
)

// Промежуточное ПО для аутентификации
func Auth(next http.HandlerFunc, env config.Environment) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(env.TodoPassword) > 0 {

			cookie, err := r.Cookie("token")

			if err != nil {
				setBadAuthentification(w, http.StatusUnauthorized, err.Error())
				return
			}

			if len(cookie.Value) == 0 {
				setBadAuthentification(w, http.StatusUnauthorized, "token is empty")
				return
			}

			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				return []byte(env.TodoPassword), nil
			})

			if err != nil || !token.Valid {
				setBadAuthentification(w, http.StatusUnauthorized, "not a valid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				setBadAuthentification(w, http.StatusUnauthorized, "authentification required")
				return
			}

			passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(env.TodoPassword)))

			if claims["hash"] != passHash {
				setBadAuthentification(w, http.StatusUnauthorized, "password has changed")
				return
			}
		}

		next(w, r)
	})
}

// Получить jwt токен
func GetJwt(password string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
	})

	strToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return strToken, nil
}

// Отправляет в ответ ошибку в о то что аутентификация не удалась
func setBadAuthentification(w http.ResponseWriter, statusCode int, err string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	errEncode := json.NewEncoder(w).Encode(err)

	if errEncode != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
	}
}
