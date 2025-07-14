package storage

import "database/sql"

// TODO интерфейс storage и реализация

type Storage struct {
	Db *sql.DB
}

type ManagerDB interface {
	AddUser(user User)
	DelUser(id int)
	Init() error
}
