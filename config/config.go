package config

import (
	"log"

	"github.com/spf13/viper"
)

var PORT string

func init() {
	viper.SetConfigName("iitkbucks-config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[WARNING] Unable to locate configuration file")
	}

	viper.AutomaticEnv()

	PORT = viper.GetString("port")

}
