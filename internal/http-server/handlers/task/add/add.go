package add

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
}

type Request struct {
	Uuid string `json:"uuid" validate:"required"`
	Task string `json:"task" validate:"required"`
}

type TASKAdder interface {
	AddTask(uuid string, taskToSave string) error
}

func New(log *slog.Logger, taskSave TASKAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.add.New"

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
		tasks := req.Task

		err = taskSave.AddTask(uuid, tasks)
		if err != nil {
			log.Info("failed to add task", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add task"))

			return
		}

		log.Info("task added")

		render.JSON(w, r, Response{
			Response: response.OK(),
		})

		return
	}
}
