package handler

import (
	"net/http"
	"strings"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	longURL, err := h.Datastore.GetLongURL(shortCode)
	if err != nil {
		h.Logger.Error().Err(err).Str("path", shortCode).Msg("Handler.GetURL unable to get the long url")
		http.Error(w, "Unable to get the long url", http.StatusInternalServerError)
		return
	}

	if longURL == "" {
		h.Logger.Warn().Err(err).Str("short_code", shortCode).Msg("Handler.GetURL 404")
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
