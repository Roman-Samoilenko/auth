package main

import (
	"auth/internal/handlers"
	"auth/internal/storage"
	"auth/internal/storage/psql"
	"auth/pkg/middleware"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки .env", "ошибка", err)
	}

	router := http.NewServeMux()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	addr := ":" + os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr: addr,
		Handler: middleware.RecoverMiddleware(
			middleware.Logging(router)),
	}

	DBtype := os.Getenv("DB_TYPE")
	var mdb storage.ManagerDB
	switch DBtype {
	case "postgres":
		mdb = &psql.Psql{}
		if err := mdb.Init(); err != nil {
			slog.Error("ошибка подключения к БД", err)
			return
		}
	default:
		mdb = &psql.Psql{}
		if err := mdb.Init(); err != nil {
			slog.Error("ошибка подключения к БД", err)
			return
		}
	}

	router.HandleFunc("/auth", handlers.Auth)

	go server.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	sign := <-stop

	slog.Info("Остановка graceful shutdown, сигнал:", sign)
}
