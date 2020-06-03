package conf

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func GetConfig() *viper.Viper {
	return conf
}

func LoadConfig() {
	log.Debug("Reading config file")
	viper.SetConfigName("bot-config")
	viper.AddConfigPath("$HOME/.telegram-rss-bot")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./.telegram-rss-bot")
	err := viper.ReadInConfig()

	if err != nil {
		log.Panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf = viper.GetViper()
	setLoggerLevel()
}

func setLoggerLevel() {
	switch conf.GetString("log_level") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}
}
