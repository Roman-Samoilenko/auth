package middleware

import (
	"log/slog"
	"net/http"
)

// Recover перехватывает панику, предотвращает падение всего сервиса
func Recover(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Перехвачена паника", "ошибка", err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
