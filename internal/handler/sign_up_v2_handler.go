package handler

import (
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) SignUpV2Handler(w http.ResponseWriter, r *http.Request) {
	var in model.SignUpV2Data
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	tokenFromDB, err := h.uc.GetTokenByName(r.Context(), "recipe")
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	if token != tokenFromDB {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	result, err := h.uc.SignUpV2(r.Context(), in)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("sign up: %s", string(res))

	_, err = w.Write(res)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
