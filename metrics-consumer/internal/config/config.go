package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Configs struct {
	Kafka Kafka `mapstructure:"kafka"`
}

type (
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	}
)

var (
	once     sync.Once
	instance *Configs
)

func GetConfig() *Configs {
	once.Do(func() {
		pwd, err := os.Getwd()
		if err != nil {
			log.Println("err in config Getwd func: ", err)
			os.Exit(1)
		}

		pathConfig := pwd + "/config.yml"

		viper.SetConfigFile(pathConfig)

		err = viper.ReadInConfig()
		if err != nil {
			log.Println("err in load config:", err)
			os.Exit(1)
		}

		instance = &Configs{}
		if err = viper.Unmarshal(instance); err != nil {
			log.Println("err in marshal load config:", err)
			os.Exit(1)
		}
	})

	return instance
}
