package main

import (
	"fmt"
	"github.com/yahorchik/TaskStorage/internal/config"
	"github.com/yahorchik/TaskStorage/internal/lib/sl"
	"github.com/yahorchik/TaskStorage/internal/logger"
	"github.com/yahorchik/TaskStorage/internal/storage/sqlite"
	"os"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Task Storage will be started")
	log.Debug("Debug mode ON")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init Storage", sl.Err(err))
		os.Exit(1)
	}
	test := sqlite.Task{
		Id:         0,
		Name:       "test5",
		Desk:       "chlen",
		Tags:       "pizda",
		CreateData: time.Now(),
		Deadline:   time.Now(),
	}
	st := storage

	id, err := st.SaveTask(test)
	if err != nil {
		log.Error("test", sl.Err(err))
		os.Exit(12)
	}

	fmt.Println(id)
}
