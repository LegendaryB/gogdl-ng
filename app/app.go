package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gogdl-ng/gogdl-ng/app/api/v1"
	"github.com/gogdl-ng/gogdl-ng/app/config"
	"github.com/gogdl-ng/gogdl-ng/app/download"
	"github.com/gogdl-ng/gogdl-ng/app/gdrive"
	"github.com/gogdl-ng/gogdl-ng/app/logging"
	"github.com/gorilla/mux"
)

func Run() {
	conf, err := config.NewConfigurationFromFile()

	if err != nil {
		log.Fatalf("Failed to retrieve app configuration. %v", err)
	}

	logger, err := logging.NewLogger(conf.Application.LogFilePath)

	if err != nil {
		log.Fatalf("Failed to initialize logger. %s", err)
	}

	drive, err := gdrive.NewDriveService(conf, logger)

	if err != nil {
		logger.Fatalf("Failed to create Google Drive service. %v", err)
	}

	jobManager, err := download.NewJobManager(logger, conf, drive)

	if err != nil {
		logger.Fatalf("Failed to create job manager. %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	controller := api.NewJobController(logger, jobManager)

	router.HandleFunc("/jobs", controller.CreateJob()).Methods("POST")

	go listenAndServe(router, conf.Application.ListenPort)

	jobManager.Run()
}

func listenAndServe(router *mux.Router, listenPort int) {
	addr := fmt.Sprintf(":%d", listenPort)

	log.Fatal(http.ListenAndServe(addr, router))
}
