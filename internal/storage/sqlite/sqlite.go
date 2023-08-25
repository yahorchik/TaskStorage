package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yahorchik/TaskStorage/internal/storage"
)

type Storage struct {
	db *sql.DB
}
type Task struct {
	Name       string
	Desk       string
	Tag        string
	CreateData string
	Deadline   string
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
	   tag TEXT NOT NULL,
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

	stmt, err := s.db.Prepare("INSERT INTO task(name, desk, tag, create_data, deadline) VALUES ( ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(taskToSave.Name, taskToSave.Desk, taskToSave.Tag, taskToSave.CreateData, taskToSave.Deadline)
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

func (s *Storage) GetTask(idTask int64) (Task, error) {
	const op = "storage.sqlite.GetTask"
	stmt, err := s.db.Prepare("SELECT * FROM task WHERE id =?")
	if err != nil {
		return Task{}, fmt.Errorf("%s: %w", op, err)
	}
	var resTask Task
	var resId int
	err = stmt.QueryRow(idTask).Scan(&resId, &resTask.Name, &resTask.Desk, &resTask.Tag, &resTask.CreateData, &resTask.Deadline)
	if errors.Is(err, sql.ErrNoRows) {
		return resTask, storage.ErrTaskExists
	}
	if err != nil {
		return resTask, fmt.Errorf("%s: %w", op, err)
	}
	return resTask, nil

}
