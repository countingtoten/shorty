package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/countingtoten/shorty"
)

type NewShortURLRequest struct {
	UserID  shorty.UserID  `json:"user_id"`
	LongURL shorty.LongURL `json:"url"`
}

type NewShortURLResponse struct {
	ShortURL shorty.ShortURL `json:"short_url"`
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	_, err := io.CopyN(buf, r.Body, r.ContentLength)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Handler.CreateShortURL unable to read request body")
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	req := &NewShortURLRequest{}
	err = json.Unmarshal(buf.Bytes(), req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Handler.CreateShortURL unable to parse request")
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.Datastore.CreateShortURL(req.UserID, req.LongURL)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Handler.CreateShortURL unable to create short url")
		http.Error(w, "Unable to create short url", http.StatusInternalServerError)
		return
	}

	resp := &NewShortURLResponse{
		ShortURL: shortURL,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
