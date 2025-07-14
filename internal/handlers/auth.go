package handlers

import (
	"auth/internal/storage"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запросы", http.StatusBadRequest)
		return
	}
	var user storage.User
	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("Ошибка при получении данных с фронтенда", "ошибка", err)
		return
	}
	//TODO: доделать логику (уникальные логины, отправка логина и хэша в БД)

	slog.Info("Получен JSON файл", "user", user)

	resp := map[string]string{
		"message": "Регистрация успешна",
		"login":   user.Login,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка создания ответа", http.StatusInternalServerError)
	}
}
