package middleware

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Внутреняя ошибка", "ошибка", err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
