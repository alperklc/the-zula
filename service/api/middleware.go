package api

import (
	"context"
	"net/http"
	"os"

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

func GetAuthMiddleware(log zerolog.Logger, domain, key string) func(_ http.Handler) http.Handler {
	ctx := context.Background()

	// Initiate the authorization by providing a zitadel configuration and a verifier.
	// This example will use OAuth2 Introspection for this, therefore you will also need to provide the downloaded api key.json
	authZ, err := authorization.New(ctx, zitadel.New(domain), oauth.DefaultAuthorization(key))
	if err != nil {
		log.Error().Msgf("zitadel sdk could not initialize %s", err)
		os.Exit(1)
	}

	return middleware.New(authZ).RequireAuthorization()
	/* return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(res, req)
		})
	} */
}
