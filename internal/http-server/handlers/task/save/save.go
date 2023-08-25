package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/yahorchik/TaskStorage/internal/lib/api/response"
	"github.com/yahorchik/TaskStorage/internal/lib/logger/sl"
	"github.com/yahorchik/TaskStorage/internal/storage"
	"github.com/yahorchik/TaskStorage/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

//type Task struct {
//	Name       string
//	Desk       string
//	Tag        string
//	CreateData time.Time
//	Deadline   time.Time
//}

type Request struct {
	Name       string `json:"name" validate:"required, name"`
	Desk       string `json:"desk"`
	Tag        string `json:"tag,omitempty"`
	CreateData string `json:"createdata"`
	Deadline   string `json:"deadline"`
}

type Response struct {
	resp.Response
	Id int64
}

type TaskSaver interface {
	SaveTask(task sqlite.Task) (sqlite.Task, error)
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		sreq := sqlite.Task{
			Name:       req.Name,
			Desk:       req.Desk,
			Tag:        req.Tag,
			CreateData: req.CreateData,
			Deadline:   req.Deadline,
		}

		id, err := taskSaver.SaveTask(sreq)
		if errors.Is(err, storage.ErrTaskExists) {
			log.Info("task already exists", slog.String("task", req.Name))

			render.JSON(w, r, resp.Error("task already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add task"))

			return
		}

	}

}
