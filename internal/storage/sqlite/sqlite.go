package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yahorchik/TaskStorage/internal/storage"
	"time"
)

type Storage struct {
	db *sql.DB
}
type Task struct {
	Id         int64
	Name       string
	Desk       string
	Tags       string
	CreateData time.Time
	Deadline   time.Time
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS task(
	   id INTEGER PRIMARY KEY,
	   name TEXT NOT NULL UNIQUE,
	   desk TEXT NOT NULL,
	   tags TEXT NOT NULL,
	   create_data TEXT NOT NULL,
	   deadline TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_id on task(name)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveTask(taskToSave Task) (int64, error) {
	const op = "storage.sqlite.SaveTask"

	stmt, err := s.db.Prepare("INSERT INTO task(name, desk, tags, create_data, deadline) VALUES ( ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(taskToSave.Name, taskToSave.Desk, taskToSave.Tags, taskToSave.CreateData, taskToSave.Deadline)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintCheck {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrTaskExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
