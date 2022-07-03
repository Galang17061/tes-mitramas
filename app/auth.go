package app

import (
	"context"
	"net/http"
	"strings"

	"tes-mitramas/controllers"
	u "tes-mitramas/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/login", "/api/register"}
		requestPath := r.URL.Path

		for _, value := range notAuth {

			if strings.Contains(requestPath, value) {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(false, "Auth token tidak ditemukan")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Token tidak valid")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1]

		type MyCustomClaims struct {
			Id       int64  `json:"id"`
			Username string `json:"username"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(tokenPart, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("M1TR4MAS"), nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					response = u.Message(false, "Token tidak valid")
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					u.Respond(w, response)
					return
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					response = u.Message(false, "Token expired")
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					u.Respond(w, response)
					return
				} else {
					response = u.Message(false, "Autentikasi token gagal "+err.Error())
					w.WriteHeader(http.StatusForbidden)
					w.Header().Add("Content-Type", "application/json")
					u.Respond(w, response)
					return
				}
			}
		}

		claims := &MyCustomClaims{}
		if claim, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			claims = claim
		} else {
			response = u.Message(false, "Token tidak valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		cont := context.WithValue(r.Context(), "values", controllers.Valcontex{Id: claims.Id, Username: claims.Username})
		r = r.WithContext(cont)

		next.ServeHTTP(w, r)

	})
}
