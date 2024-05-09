package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/alperklc/the-zula/service/api"
	notesCtrl "github.com/alperklc/the-zula/service/controller/notes"
	notesReferencesCtrl "github.com/alperklc/the-zula/service/controller/notesReferences"
	"github.com/alperklc/the-zula/service/infrastructure/db"
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"
	"github.com/alperklc/the-zula/service/infrastructure/environment"
	"github.com/alperklc/the-zula/service/infrastructure/logger"
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

	nr := notes.NewDb(d)
	ndr := notesDrafts.NewDb(d)
	nrr := notesReferences.NewDb(d)

	noteController := notesCtrl.NewNotesController(nr, ndr)
	notesReferencesController := notesReferencesCtrl.NewNotesReferencesController(nr, nrr)

	a := api.NewApi(noteController, notesReferencesController)
	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))

	authMW := api.GetAuthMiddleware(l, config.AuthDomain, config.AuthKeyFilePath)
	r.Use(authMW)

	api.HandlerFromMux(a, r)

	server := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", config.Port)),
	}

	l.Info().Msgf("the zula service started on http://localhost:%d/", config.Port)
	log.Fatal(server.ListenAndServe())
}
