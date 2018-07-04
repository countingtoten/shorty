package handler

import (
	"net/http"

	"github.com/countingtoten/shorty"
	"github.com/rs/zerolog"
)

type Handler struct {
	shorty.Datastore
	zerolog.Logger
}

func New(ds shorty.Datastore, logger zerolog.Logger) *Handler {
	return &Handler{
		Datastore: ds,
		Logger:    logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/new":
		if r.Method == http.MethodPost {
			h.CreateShortURL(w, r)
		} else {
			http.NotFound(w, r)
		}
	default:
		h.GetURL(w, r)
	}
}
