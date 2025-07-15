package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Log представляет собой структуру для логирования с JSON-тегами
type Log struct {
	Message    string `json:"message"`
	Time       string `json:"time"`
	UserID     string `json:"user_id,omitempty"`
	Provider   string `json:"provider,omitempty"` // "jwt" или "github"
	Method     string `json:"method,omitempty"`   // HTTP метод
	Path       string `json:"path,omitempty"`     // URL путь
	StatusCode int    `json:"status_code,omitempty"`
	Error      string `json:"error,omitempty"`
}

func Logging(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		logEntry := Log{
			Message: r.RequestURI,
			Time:    time.Now().UTC().Format(time.RFC3339),
			Method:  r.Method,
			Path:    r.URL.Path,
		}

		slog.Info("HTTP запрос",
			"метод", logEntry.Method,
			"путь", logEntry.Path,
		)

	})

}
