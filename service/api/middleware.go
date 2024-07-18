package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/rs/zerolog"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"

	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type APIErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetLoggingMiddleware(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()

			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("remote_addr", req.RemoteAddr).
				Msg("Request received")

			next.ServeHTTP(res, req)

			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("remote_addr", req.RemoteAddr).
				Dur("duration", time.Since(start)).
				Msg("Request completed")
		})
	}
}

func GetAuthorizationMiddleware(log zerolog.Logger, domain, key string) func(_ http.Handler) http.Handler {
	ctx := context.Background()

	zm := zitadel.New(domain)
	authZ, err := authorization.New(ctx, zm, oauth.DefaultAuthorization(key))

	if err != nil {
		log.Error().Msgf("zitadel sdk could not initialize %s", err)
		os.Exit(1)
	}

	return middleware.New(authZ).RequireAuthorization()
}

func GetAuthenticationMiddleware(notes notes.Collection, domain, key string) func(_ http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			user := authorization.UserID(req.Context())

			fmt.Println(user)
			fmt.Println(req.URL.Path)

			next.ServeHTTP(res, req) // req.WithContext(ctx))
		})
	}
}
