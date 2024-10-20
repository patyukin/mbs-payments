package handler

import (
	"auth-telegram/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) SignInVerifyHandler(w http.ResponseWriter, r *http.Request) {
	var signInVerifyData model.SignInVerifyData

	if err := json.NewDecoder(r.Body).Decode(&signInVerifyData); err != nil {
		log.Error().Err(err).Msgf("failed to decode sign in verify data, error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.SignInVerify(r.Context(), signInVerifyData.Code)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in verify, error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		log.Error().Err(err).Msgf("failed to encode tokens, error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
