package config

import (
	"log"

	"github.com/spf13/viper"
)

var PORT string
var MY_URL string

var MAX_PEERS int
var POTENTIAL_PEERS []string

func init() {
	viper.SetConfigName("iitkbucks-config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[WARNING] Unable to locate configuration file")
	}

	viper.AutomaticEnv()

	PORT = viper.GetString("port")
	MY_URL = viper.GetString("myUrl")

	MAX_PEERS = viper.GetInt("maxPeers")
	POTENTIAL_PEERS = viper.GetStringSlice("potentialPeers")
}
