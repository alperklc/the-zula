package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alperklc/the-zula/service/api"
	"github.com/alperklc/the-zula/service/infrastructure/auth"
	"github.com/alperklc/the-zula/service/infrastructure/cache"
	"github.com/alperklc/the-zula/service/infrastructure/db"
	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesChanges"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"
	"github.com/alperklc/the-zula/service/infrastructure/db/pageContent"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"
	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	"github.com/alperklc/the-zula/service/infrastructure/environment"
	"github.com/alperklc/the-zula/service/infrastructure/logger"
	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	"github.com/alperklc/the-zula/service/infrastructure/webScraper"
	bookmarksService "github.com/alperklc/the-zula/service/services/bookmarks"
	importerService "github.com/alperklc/the-zula/service/services/importer"
	notesService "github.com/alperklc/the-zula/service/services/notes"
	referencesService "github.com/alperklc/the-zula/service/services/references"
	userActivityService "github.com/alperklc/the-zula/service/services/userActivity"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	config := environment.Read()
	logger.Init(config.LogLevel)

	l := logger.Get()
	d := db.Connect(config.MongoURI)
	mq, errMq := messagequeue.New(config.RabbitMqUri)
	if errMq != nil {
		l.Fatal().Msgf("Error connecting to the rabbitmq: %s", errMq)
	}
	mqp := mqpublisher.New(mq.Channel)

	swagger, err := api.GetSwagger()
	if err != nil {
		l.Fatal().Msgf("Error loading swagger spec\n: %s", err)
	}

	swagger.Servers = nil

	ac := auth.NewAuthClient(fmt.Sprintf("https://%s/", config.AuthDomain), config.AuthServiceAccountUser, config.AuthServiceAccountSecret)

	nr := notes.NewDb(d)
	ncr := notesChanges.NewDb(d)
	ndr := notesDrafts.NewDb(d)
	nrr := references.NewDb(d)
	uad := useractivity.NewDb(d)
	b := bookmarks.NewDb(d)
	pc := pageContent.NewDb(d)

	wsc := webScraper.NewWebScraper()

	ums, errMemstore := cache.NewCache[usersService.User](1 * time.Hour)
	if errMemstore != nil {
		l.Fatal().Msgf("Error memory store: %s", errMemstore)
	}
	uams, errMemstore2 := cache.NewCache[[]useractivity.ActivityGraphEntry](1 * time.Hour)
	if errMemstore2 != nil {
		l.Fatal().Msgf("Error memory store: %s", errMemstore2)
	}
	usms, errMemstore3 := cache.NewCache[[]useractivity.UsageStatisticsEntry](1 * time.Hour)
	if errMemstore3 != nil {
		l.Fatal().Msgf("Error memory store: %s", errMemstore3)
	}

	us := usersService.NewService(ac, ums, mqp)
	bs := bookmarksService.NewService(l, us, b, pc, wsc, mqp)
	nrs := referencesService.NewService(nr, nrr)
	ns := notesService.NewService(us, nr, ncr, ndr, nrs, mqp)
	uas := userActivityService.NewService(uams, usms, us, uad, ns, bs)
	is := importerService.NewService(nr, nrr, b, pc, uad)

	uasmq := userActivityService.NewMqConsumer(l, uad, mq)
	uasmq.Start()

	rsmq := referencesService.NewMqConsumer(l, nrs, mq)
	rsmq.Start()

	usmq := usersService.NewMqConsumer(l, us, mq)
	usmq.Start()

	hub := api.NewHub()
	go hub.Run()

	ws := api.NewNotifier(l, mq, *hub)
	ws.SendNotification()

	a := api.NewApi(us, uas, bs, ns, is, *hub)
	r := chi.NewRouter()

	r.Use(api.GetLoggingMiddleware(l))
	r.Use(api.GetAuthorizationMiddleware(l, config.AuthDomain, config.AuthKeyFilePath))
	r.Use(api.GetAuthenticationMiddleware(nr, config.AuthDomain, config.AuthKeyFilePath))

	staticFileServer := http.FileServer(http.FS(staticFiles))

	r.Group(func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path

			_, err := staticFiles.Open(filepath.Join("static", filePath))
			if err == nil {
				staticFileServer.ServeHTTP(w, r)
			} else {
				http.ServeFile(w, r, "static/index.html")
			}
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.OapiRequestValidator(swagger))
		api.HandlerFromMux(a, r)
	})

	server := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", config.Port)),
	}
	l.Info().Msgf("the zula service started on http://localhost:%d/", config.Port)
	log.Fatal(server.ListenAndServe())
}
