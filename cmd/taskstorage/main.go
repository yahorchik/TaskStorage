package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yahorchik/TaskStorage/internal/config"
	mwLogger "github.com/yahorchik/TaskStorage/internal/http-server/middleware/logger"
	"github.com/yahorchik/TaskStorage/internal/lib/logger/sl"
	"github.com/yahorchik/TaskStorage/internal/logger"
	"github.com/yahorchik/TaskStorage/internal/storage/sqlite"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Task Storage will be started")
	log.Debug("Debug mode ON")
	log.Error("error message are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init Storage", sl.Err(err))
		os.Exit(1)
	}
	//test := sqlite.Task{
	//	Name:       "test",
	//	Desk:       "chlen",
	//	Tag:        "pizda",
	//	CreateData: time.Now(),
	//	Deadline:   time.Now(),
	//}
	_ = storage
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	//middleware
	//id, err := st.SaveTask(test)
	//if err != nil {
	//	log.Error("test", sl.Err(err))
	//	os.Exit(12)
	//}
	//test, err = st.GetTask(1)
	//if err != err {
	//	os.Exit(1)
	//}
	//fmt.Println(test.Name)
	//	fmt.Println(id)
}
