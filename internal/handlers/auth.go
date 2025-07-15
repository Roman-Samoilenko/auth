package handlers

import (
	"auth/internal/storage"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// AuthHandler структура-обёртка над storage.ManagerDB, реализующая ServeHTTP,
// из-за чего может использоваться в методе Handle, нужна для доступа к БД в обработчиках.
type AuthHandler struct {
	Mdb storage.ManagerDB
}

// ServeHTTP метод AuthHandler, для реализации интерфейса Handle
func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	a.Mdb.AddUser(user, context.Background())

	resp := map[string]string{
		"message": "Регистрация успешна",
		"login":   user.Login,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка создания ответа", http.StatusInternalServerError)
	}
}
