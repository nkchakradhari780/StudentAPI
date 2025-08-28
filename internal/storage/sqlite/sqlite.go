package sqlite

import (
	"database/sql"

	"github.com/nkchakradhari780/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3",cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_ , err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLite{Db: db}, nil
}

func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error){
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}