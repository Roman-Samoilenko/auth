package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// TODO интерфейс storage и реализация

// storage структура, которая встраивается в конкретные реализации ManagerDB,
// для их доступа к методам интерфейса managerStorage
type storage struct {
	db *sql.DB
}

// managerStorage интерфейс для установки и закрытия соединения с БД
type managerStorage interface {
	Init(DBType string) error
	Close(DBType string) error
}

// managerTable интерфейс работы с таблицей пользователей
type managerTable interface {
	AddUser(user User, ctx context.Context) error
	DelUser(login string, ctx context.Context) error
}

// ManagerDB представляет основной интерфейс для слоя данных.
// Объединяет функциональность managerStorage и managerTable.
// Используется в API для абстракции от конкретной реализации БД, обеспечивая независимость
// бизнес-логики от деталей хранения данных и позволяя легко менять тип БД.
type ManagerDB interface {
	managerStorage
	managerTable
}

// newStorage конструктор для storage
func newStorage(db *sql.DB) storage {
	return storage{db: db}
}

// Init открывает соединение с БД
func (s *storage) Init(DBT string) error {
	pass := os.Getenv(DBT + "_PASS")
	host := os.Getenv(DBT + "_HOST")
	port := os.Getenv(DBT + "_PORT")
	user := os.Getenv(DBT + "_USER")
	dbname := os.Getenv(DBT + "_NAME")

	reqInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", reqInfo)
	if err != nil {
		return fmt.Errorf("ошибка подключения к %S: %w", DBT, err)
	}

	//TODO: обернуть все ошибки и всегда возвращать обёрнутые ошибки

	if err = db.Ping(); err != nil {
		return fmt.Errorf("БД %S не пингуется %w", DBT, err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		pass_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if SQLRequest, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %w, SQL запрос: %s", err, SQLRequest)
	}

	s.db = db

	return nil
}

// Close закрывает соединение с БД
func (s *storage) Close(DBT string) error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("ошибка: %w, при закрытии соединения с БД: %s", err, DBT)
	}
	return nil
}
