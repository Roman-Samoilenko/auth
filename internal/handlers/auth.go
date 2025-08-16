package handlers

import (
	"auth/internal/service"
	"auth/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// AuthHandler структура-обёртка над storage.ManagerDB, реализующая ServeHTTP,
// из-за чего может использоваться в методе Handle, нужна для доступа к БД в обработчиках.
type AuthHandler struct {
	mt storage.ManagerTable
}

func NewAuthHandler(mt storage.ManagerTable) *AuthHandler {
	return &AuthHandler{mt: mt}
}

// ServeHTTP метод AuthHandler, для реализации интерфейса Handle
func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(r *http.Request) {
		r.Body.Close()
	}(r)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		msg := "Ошибка чтения тела запросы"
		http.Error(w, msg, http.StatusBadRequest)
		slog.Error(msg,
			"ошибка", err,
			"логин", "",
			"handler", "auth")
		reply(w, "Ошибка чтения тела запросы", "", "")
		return
	}
	var user storage.User
	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("Ошибка при получении данных пользователя",
			"ошибка", err,
			"логин", user.Login,
			"handler", "auth")
		reply(w, "Ошибка при получении данных пользователя", "", "")
		return
	}

	exist, err := a.mt.LoginExists(ctx, user.Login)
	if err != nil {
		slog.Error("ошибка при проверки сущ. пользователя",
			"ошибка", err,
			"логин", user.Login,
			"handler", "auth")
		reply(w, fmt.Sprintf("ошибка: %w при проверки сущ. пользователя: %s", err, user.Login), user.Login, "")
		return
	}
	if !exist {
		if err := a.mt.AddUser(ctx, user); err != nil {
			slog.Error("ошибка добавления логина",
				"ошибка", err,
				"логин", user.Login,
				"handler", "auth")
			reply(w, fmt.Sprintf("ошибка: %w добавления логина: %s", err, user.Login), user.Login, "")
			return
		}

		strToken, err := service.CreateJWt(user.Login)
		if err != nil {
			slog.Error("ошибка подписания JWT токена",
				"ошибка", err,
				"логин", user.Login,
				"handler", "auth")
			reply(w, fmt.Sprintf("ошибка подписания JWT токена: %w", err), user.Login, strToken)
			return
		}
		reply(w, "Регистрация успешна", user.Login, strToken)

	} else {
		correct, err := a.mt.CheckPassHash(ctx, user.Pass, user.Login)
		if err != nil {
			slog.Error("ошибка проверка пароля",
				"ошибка", err,
				"логин", user.Login,
				"handler", "auth")
			reply(w, fmt.Sprintf("ошибка проверка пароля: %w", err), user.Login, "")
			return
		}
		if !correct {
			reply(w, "неправильный пароль для этого логина", user.Login, "")
			return
		}
		strToken, err := service.CreateJWt(user.Login)
		if err != nil {
			slog.Error("ошибка подписания JWT токена",
				"ошибка", err,
				"логин", user.Login,
				"handler", "auth")
			reply(w, fmt.Sprintf("ошибка подписания JWT токена: %w", err), user.Login, strToken)
			return
		}
		reply(w, "вход успешен", user.Login, strToken)
	}
	return
}

func reply(w http.ResponseWriter, msg string, login string, strToken string) {
	resp := map[string]string{
		"message": msg,
		"login":   login,
		"token":   strToken,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    strToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   3600,
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка создания ответа", http.StatusInternalServerError)
	}
}
