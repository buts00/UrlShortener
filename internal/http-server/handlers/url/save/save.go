package save

import (
	"errors"
	"fmt"
	resp "github.com/buts00/UrlShortener/internal/lib/app/response"
	"github.com/buts00/UrlShortener/internal/lib/random"
	"github.com/buts00/UrlShortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 6

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type urlSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSave urlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.url.save.New"
		log.With(
			slog.String("op", op),
			slog.String("id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		fmt.Println(r.Body)
		if err != nil {
			log.Error("failed to decode request body ", err)
			render.JSON(w, r, resp.Error("failed to decode request "))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("inappropriate data load ", err)
			render.JSON(w, r, resp.Error("inappropriate data load "))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomAlias(aliasLength)
		}
		id, err := urlSave.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("failed to add url ", err)
			render.JSON(w, r, resp.Error("failed to add url "))
			return
		}

		log.Info("url added", slog.String("url", req.URL), slog.Int64("id", id))
		render.JSON(w, r, resp.Ok())

	}
}
