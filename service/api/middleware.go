package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type APIErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		user, ok := req.Header["Remote-User"]
		if !ok {
			response := APIErrorResponse{Message: "User header is missing", Status: "NOT_AUTHENTICATED"}

			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(response)
			return
		}

		ctxWithUser := context.WithValue(req.Context(), "user", user[0])
		next.ServeHTTP(res, req.WithContext(ctxWithUser))
		return
	})
}
