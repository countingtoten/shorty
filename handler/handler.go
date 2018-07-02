package handler

import (
	"net/http"

	"github.com/countingtoten/shorty"
	"github.com/rs/zerolog"
)

type Handler struct {
	UserData  map[shorty.UserID]*shorty.User
	ShortURLs map[shorty.ShortURL]*shorty.URL
	zerolog.Logger
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
