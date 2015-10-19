package common

import (
	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	BuildDate  string
	CommitHash string
	Port       int
	Host       string
)

func ReadEnvironment() {
	// override configuration with environment vars
	// example: GREEDY_PORT=/tmp/greedy.sqlite
	viper.SetEnvPrefix("PODDY")
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "0.0.0.0")
	viper.AutomaticEnv()

	Port = viper.GetInt("port")
	Host = viper.GetString("host")

	log.Info("environment", "host", Host, "port", Port)
	log.Info("meta-info", "builddate", BuildDate, "commithash", CommitHash)
}
