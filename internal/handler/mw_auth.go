package handler

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessToken, err := GetBearerToken(r)
		if err != nil {
			h.HandleError(w, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return h.uc.GetJWTToken(), nil
		})

		if err != nil || !token.Valid {
			h.HandleError(w, http.StatusUnauthorized, err.Error())
			return
		}

		id := token.Claims.(jwt.MapClaims)["id"].(string)

		r.Header.Set(HeaderUserID, id)

		next.ServeHTTP(w, r)
	})
}
