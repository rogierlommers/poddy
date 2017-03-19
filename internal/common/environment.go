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
	Storage    string
	Watch      string
	Self       string
)

func ReadEnvironment() {
	// override configuration with environment vars
	// example: PODDY_PORT=8080
	viper.SetEnvPrefix("PODDY")
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("storage", "./storage")
	viper.SetDefault("watch", "./watch")
	viper.SetDefault("self", "http://poddy.lommers.org")
	viper.AutomaticEnv()

	Port = viper.GetInt("port")
	Host = viper.GetString("host")
	Storage = viper.GetString("storage")
	Watch = viper.GetString("watch")
	Self = viper.GetString("self")

	log.Info("environment info loaded", "host", Host, "port", Port, "storage", Storage, "watch", Watch)
	log.Info("build information", "builddate", BuildDate, "commithash", CommitHash)
}
