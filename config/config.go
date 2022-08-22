package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func init() {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("err read config" + err.Error())
	}

	config = v
}

func GetConfig() *viper.Viper {
	return config
}
