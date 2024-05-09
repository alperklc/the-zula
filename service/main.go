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

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	DB *mongo.Database
}

func main() {
	config := environment.Read()
	logger.Init(config.LogLevel)

	l := logger.Get()
	d := db.Connect(config.MongoURI)

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
	}

	swagger.Servers = nil

	notesRepository := notes.NewNotesRepository(d)
	notesReferencesRepository := notesReferences.NewRepository(d)
	notesDraftsRepository := notesDrafts.NewNotesDraftsRepository(d)

	noteController := notesCtrl.NewNotesController(notesRepository, notesDraftsRepository)
	notesReferencesController := notesReferencesCtrl.NewNotesReferencesController(notesRepository, notesReferencesRepository)

	a := api.NewApi(noteController, notesReferencesController)
	r := mux.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(api.AuthMiddleware)
	api.HandlerFromMux(a, r)

	server := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", config.Port)),
	}

	l.Info().Msgf("the zula service started on http://localhost:%d/", config.Port)
	log.Fatal(server.ListenAndServe())
}
