package handler

import (
	"auth-telegram/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) SignInV2Handler(w http.ResponseWriter, r *http.Request) {
	var signInData model.SignInData
	if err := json.NewDecoder(r.Body).Decode(&signInData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.SignInV2(r.Context(), signInData)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		log.Error().Err(err).Msgf("failed to encode tokens, error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
