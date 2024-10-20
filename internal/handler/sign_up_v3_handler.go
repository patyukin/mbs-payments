package handler

import (
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (h *Handler) SignUpV3Handler(w http.ResponseWriter, r *http.Request) {
	var in model.SignUpV2Data
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug().Msgf("data in: %v", in)

	tokenWithService := r.Header.Get("Authorization")
	parts := strings.SplitN(tokenWithService, ":", 2)

	if len(parts) != 2 {
		log.Error().Msgf("failed token: %s", tokenWithService)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	service := parts[0]
	token := parts[1]

	tokenFromDB, err := h.uc.GetTokenByName(r.Context(), service)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up in h.uc.GetTokenByName, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	log.Info().Msgf("service: %s, token: %s, tokenFromDB: %s", service, token, tokenFromDB)
	if token != tokenFromDB {
		log.Error().Err(err).Msgf("failed to sign up tokens do not match, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	log.Debug().Msgf("data in: %v", in)

	result, err := h.uc.SignUpV2(r.Context(), in)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up in h.uc.SignUpV2, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
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
