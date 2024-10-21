package handler

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get(HeaderUserID)
		userUUID, err := uuid.Parse(id)
		if err != nil {
			h.HandleError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userAuthInfo, err := h.uc.GetUserAuthInfoByToken(r.Context(), userUUID.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.HandleError(w, http.StatusUnauthorized, err.Error())
				return
			}

			h.HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if userAuthInfo.Role != "admin" {
			h.HandleError(w, http.StatusForbidden, "Forbidden")
			return
		}

		r.Header.Set(HeaderUserRole, userAuthInfo.Role)

		next.ServeHTTP(w, r)
	})
}
