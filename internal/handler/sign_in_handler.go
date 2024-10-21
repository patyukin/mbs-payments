package handler

import (
	"auth-telegram/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInData model.SignInData

	if err := json.NewDecoder(r.Body).Decode(&signInData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.uc.SignIn(r.Context(), signInData)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
