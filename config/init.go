package config

import (
	"log"

	"github.com/spf13/viper"
)

var Conf *viper.Viper

func InitConf(configPath string) {
	Conf = viper.New()
	Conf.SetConfigFile(configPath)
	if err := Conf.ReadInConfig(); err != nil {
		log.Fatalf("can not read config: %v\n, err: %v", configPath, err)
	}
}
