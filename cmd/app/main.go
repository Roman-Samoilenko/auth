package main

import (
	"auth/internal/handlers"
	"auth/internal/storage"
	"auth/pkg/middleware"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки .env", "ошибка", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := http.NewServeMux()
	server := &http.Server{
		Addr: ":" + os.Getenv("SERVER_PORT"),
		Handler: middleware.Recover(
			middleware.Logging(router)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	var mdb storage.ManagerDB

	DBType := os.Getenv("DB_TYPE")
	switch DBType {
	case "postgres":
		mdb = storage.NewPsql()
	default:
		mdb = storage.NewPsql()
	}
	if err := mdb.Init(DBType); err != nil {
		slog.Error("ошибка подключения к БД", err)
		return
	}

	authHandler := &handlers.AuthHandler{Mdb: mdb}
	router.Handle("POST /auth", authHandler)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Ошибка при ListenAndServe()", err)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	sign := <-stop

	err := mdb.Close(DBType)
	if err != nil {
		slog.Error("ошибка", err)
	} else {
		slog.Info("Соединение с БД успешно разорвано")
	}

	slog.Info("Остановка graceful shutdown", "сигнал", sign)
}
