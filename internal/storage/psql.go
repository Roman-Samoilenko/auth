package storage

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

// Psql реализация managerTable с БД PostgreSQL
type Psql struct {
	storage
}

// NewPsql конструктор для Psql
func NewPsql() *Psql {
	return &Psql{
		newStorage(nil),
	}
}

// AddUser метод добавления пользователя в БД PostgreSQL
func (p *Psql) AddUser(user User, ctx context.Context) error {
	//TODO: всё писать с контекстом и исправить этот бардак
	pass_hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Ошибка при хэшировании пароля", "ошибка", err)
		return fmt.Errorf(
			"ошибка при хэшировании пароля: %w",
			err,
		)
	}
	SQLadd := "INSERT INTO users (login, pass_hash) VALUES ($1, $2)"

	p.storage.db.Exec(SQLadd, user.Login, string(pass_hash))

	return err
}

// DelUser метод удаления пользователя в БД PostgreSQL
func (p *Psql) DelUser(login string, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
