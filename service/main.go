package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/alperklc/the-zula/service/api"
	"github.com/alperklc/the-zula/service/infrastructure/auth"
	"github.com/alperklc/the-zula/service/infrastructure/db"
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"
	"github.com/alperklc/the-zula/service/infrastructure/environment"
	"github.com/alperklc/the-zula/service/infrastructure/logger"
	notesService "github.com/alperklc/the-zula/service/services/notes"
	notesReferencesService "github.com/alperklc/the-zula/service/services/notesReferences"
	usersService "github.com/alperklc/the-zula/service/services/users"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	config := environment.Read()
	logger.Init(config.LogLevel)

	l := logger.Get()
	d := db.Connect(config.MongoURI)

	swagger, err := api.GetSwagger()
	if err != nil {
		l.Fatal().Msgf("Error loading swagger spec\n: %s", err)
	}

	swagger.Servers = nil

	ac := auth.NewAuthClient(fmt.Sprintf("https://%s/", config.AuthDomain), config.AuthServiceAccountUser, config.AuthServiceAccountSecret)

	nr := notes.NewDb(d)
	ndr := notesDrafts.NewDb(d)
	nrr := notesReferences.NewDb(d)

	usersService := usersService.NewService(ac)
	noteService := notesService.NewService(usersService, nr, ndr)
	notesReferencesService := notesReferencesService.NewService(nr, nrr)

	a := api.NewApi(usersService, noteService, notesReferencesService)
	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))

	r.Use(api.GetLoggingMiddleware(l))
	r.Use(api.GetAuthorizationMiddleware(l, config.AuthDomain, config.AuthKeyFilePath))
	r.Use(api.GetAuthenticationMiddleware(nr, config.AuthDomain, config.AuthKeyFilePath))

	api.HandlerFromMux(a, r)

	server := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", config.Port)),
	}

	l.Info().Msgf("the zula service started on http://localhost:%d/", config.Port)
	log.Fatal(server.ListenAndServe())
}
