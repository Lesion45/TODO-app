package get

import (
	"TODOapp/internal/lib/api/response"
	"TODOapp/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
)

type Response struct {
	response.Response
	Tasks map[string]string
}

type Request struct {
	Uuid string `json:"uuid" validate:"required"`
}

type TASKSGetter interface {
	GetTasks(uuid string) (map[string]string, error)
}

func New(log *slog.Logger, tasksGet TASKSGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get.New"

		log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = validator.New().Struct(req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidataionError(validateErr))

			return
		}

		uuid := req.Uuid

		tasks, err := tasksGet.GetTasks(uuid)
		if err != nil {
			log.Info("failed to get task", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete task"))

			return
		}

		log.Info("tasks received")

		render.JSON(w, r, Response{
			Response: response.OK(),
			Tasks:    tasks,
		})

		return
	}
}
