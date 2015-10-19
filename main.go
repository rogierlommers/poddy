package main

import (
	"github.com/rogierlommers/poddy/internal/common"
	log "gopkg.in/inconshreveable/log15.v2"
)

func main() {
	log.Info("poddy", "status", "starting service")

	// read environment vars
	common.ReadEnvironment()

}
