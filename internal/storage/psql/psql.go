package psql

import (
	"auth/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

type Psql struct {
	storage.Storage
}

func (p *Psql) Init() error {
	pass, ok := os.LookupEnv("PSQL_PASS")
	if !ok {
		slog.Error("PSQL_PASS не найдено в .env")
		return errors.New("PSQL_PASS не найдено в .env")
	}
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	user := os.Getenv("PSQL_USER")
	dbname := os.Getenv("PSQL_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		slog.Error("Ошибка подключения к postgres", "ошибка", err)
		return err
	}

	if err = db.Ping(); err != nil {
		slog.Error("БД postgres не пингуется", "ошибка", err)
		return err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		passw TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = db.Exec(createTableSQL); err != nil {
		slog.Error("Ошибка при создании таблицы", "ошибка", err)
		return err
	}

	p.Storage = storage.Storage{
		Db: db,
	}

	return nil
}

func (p *Psql) AddUser(user storage.User) {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) DelUser(id int) {
	//TODO implement me
	panic("implement me")
}
