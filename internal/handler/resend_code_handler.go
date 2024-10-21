package handler

import (
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) ResendCodeHandler(w http.ResponseWriter, r *http.Request) {
	var in model.SignUpData

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httperror.SendError(w, "invalid data", http.StatusBadRequest)
		return
	}

	if !in.IsValid() {
		log.Error().Msg("invalid input")
		httperror.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	dto, err := h.uc.ResendCode(r.Context(), in)
	if err != nil {
		log.Error().Msgf("invalid input, err: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	link := fmt.Sprintf("/start=%s", dto.Code)
	message := fmt.Sprintf("Please send this message to complete your registration: %s", link)

	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
