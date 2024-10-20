package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.Header.Get(HeaderUserID))
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := uuid.Parse(id.String())
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := h.uc.GetUserInfoByUUID(r.Context(), userUUID)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
