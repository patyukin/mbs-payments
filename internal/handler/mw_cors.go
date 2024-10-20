package handler

import (
	"net/http"
)

func (h *Handler) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Add("Access-Control-Allow-Methods", "POST")
			resp.Header().Add("Access-Control-Allow-Methods", "GET")
			resp.Header().Add("Access-Control-Allow-Headers", "Authorization")
			resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")

			return
		}
		resp.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(resp, req)
	})
}
