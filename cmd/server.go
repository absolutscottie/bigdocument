package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/absolutscottie/bigdocument/internal/config"
	"github.com/absolutscottie/bigdocument/internal/egress"
	"github.com/absolutscottie/bigdocument/internal/ingest"
	"github.com/absolutscottie/bigdocument/internal/middleware"
	"github.com/absolutscottie/bigdocument/internal/mock"
	log "github.com/sirupsen/logrus"
)

func main() {
	configLogging()

	cfg, err := config.Load(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		log.Errorf("Failed to load configuration file: %s\n", err.Error())
		return
	}

	datastore := mock.NewDatastore()
	ingest.ConfigureDatastore(datastore)
	egress.ConfigureDatastore(datastore)

	router := mux.NewRouter()
	ingest.AddHandlers(router)
	egress.AddHandlers(router)
	router.Use(middleware.LoggingMiddleware)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.BindAddr, cfg.Port),
		Handler: router,
		// realisticlly we should be setting timeouts and not using ListenAndServe() but I have limited time
	}

	server.ListenAndServe()
}

func configLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
