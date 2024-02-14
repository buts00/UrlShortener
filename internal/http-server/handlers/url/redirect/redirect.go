package redirect

import (
	"errors"
	"fmt"
	resp "github.com/buts00/UrlShortener/internal/lib/app/response"
	"github.com/buts00/UrlShortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type urlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.url.redirect.New"
		log.With(
			slog.String("op", op),
			slog.String("id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		fmt.Println(alias)
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url")
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
