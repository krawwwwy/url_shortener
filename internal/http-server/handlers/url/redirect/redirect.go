package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(

			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		urlToGet, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrNotFound) {
			msg := "no url  belong to this alias"
			log.Info(msg, slog.String("alias", alias))

			render.JSON(w, r, resp.Error(msg))

			return
		}
		if err != nil {
			msg := "failed to get url"
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, "internal server error")

			return
		}

		log.Info("successfully got url", slog.String("url", urlToGet), slog.String("alias", alias))

		http.Redirect(w, r, urlToGet, http.StatusFound)

	}
}
