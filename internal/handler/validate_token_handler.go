package handler

import (
	"auth-telegram/internal/handler/dto"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var token dto.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		log.Error().Msgf("failed to decode request: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug().Msgf("token: %s", token.Token)

	id, err := h.uc.ValidateToken(token.Token)
	if err != nil {
		log.Error().Msgf("failed to validate token: %s", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Debug().Msgf("id: %s", id)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dto.ValidateTokenResponse{UUID: id}); err != nil {
		log.Error().Msgf("failed to encode response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
