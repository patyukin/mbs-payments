package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GenerateRefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var refreshToken RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.uc.GenerateTokens(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
