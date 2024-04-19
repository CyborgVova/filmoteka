package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/filmoteka/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		login, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var hash string
		repo := postgres.NewRepository()
		repo.Conn.QueryRow(ctx, "SELECT password FROM auth WHERE login = $1", login).Scan(&hash)
		defer repo.Conn.Close(ctx)
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
