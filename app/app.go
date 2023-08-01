package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/api/v1"
	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/gorilla/mux"
)

func Run() {
	conf, err := config.NewConfigurationFromFile()

	if err != nil {
		log.Panicf("Failed to retrieve app configuration. %v", err)
	}

	logger, err := logging.NewLogger(conf.Logging)

	if err != nil {
		log.Panicf("Failed to initialize logger. %v", err)
	}

	logger.Infof("Loaded configuration file: %s", conf.GetConfigurationFolderPath())

	drive, err := gdrive.NewDriveService(conf, logger)

	if err != nil {
		logger.Fatalf("Failed to create instance of the Google Drive service. %v", err)
	}

	logger.Info("Created instance of the Google Drive service")

	jobManager, err := download.NewJobManager(logger, conf, drive)

	if err != nil {
		logger.Fatalf("failed to create instance of the Job manager. %v", err)
	}

	logger.Info("Created instance of the Job manager")

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	controller := api.NewJobController(logger, jobManager)

	router.HandleFunc("/jobs", controller.CreateJob()).Methods("POST")

	go listenAndServe(router, logger, conf.Application.ListenPort)

	logger.Info("Startup finished, continuing with running the Job manager")

	if err := jobManager.Run(); err != nil {
		logger.Fatalf("Failed to run the Job manager. %v", err)
	}
}

func listenAndServe(router *mux.Router, logger logging.Logger, listenPort int) {
	addr := fmt.Sprintf(":%d", listenPort)

	logger.Fatal(http.ListenAndServe(addr, router))
}
