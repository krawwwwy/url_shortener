package delete

import (
	"errors"
	resp "firstgomode/internal/lib/api/response"
	"firstgomode/internal/lib/logger/sl"
	"firstgomode/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrNotFound) {
			msg := "url's not exist"
			log.Info(msg, slog.String("alias", alias))

			render.JSON(w, r, resp.Error(msg))
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return

		}

		log.Info("successfully deleted url", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())

	}
}
