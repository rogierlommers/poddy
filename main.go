package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rogierlommers/poddy/internal/common"
	"github.com/rogierlommers/poddy/internal/poddy"
	log "gopkg.in/inconshreveable/log15.v2"
)

func main() {
	// read environment vars
	common.ReadEnvironment()

	// initialise mux router
	router := mux.NewRouter()

	// setup statics
	poddy.CreateStaticBox(router)

	// setup watchdirectory
	if common.Watch != "" {
		check, err := os.Stat(common.Watch)
		if err != nil || !check.IsDir() {
			log.Error("error setting up watchdirectory (point to directory?)")
		} else {
			poddy.EnableWatchdirectory(common.Watch)
		}
	}

	// http handles
	router.HandleFunc("/", poddy.IndexPage)
	router.HandleFunc("/add-podcast", poddy.AddPodcast)
	router.HandleFunc("/feed", poddy.Feed)
	router.PathPrefix("/download").Handler(http.StripPrefix("/download", http.FileServer(http.Dir("storage/"))))

	// start server
	log.Info("poddy is running/listening", "host", common.Host, "port", common.Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", common.Host, common.Port), router)
	if err != nil {
		log.Crit("daemon could not bind on interface", "host", common.Host, "port", common.Port)
		os.Exit(1)
	}
}
