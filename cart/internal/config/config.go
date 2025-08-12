package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Configs struct {
	Listen   Listen     `mapstructure:"listen"`
	Postgres DbPostgres `mapstructure:"postgres"`
	Kafka    Kafka      `mapstructure:"kafka"`
	Tracing  Tracing    `mapstructure:"tracing"`
	Metrics  Metrics    `mapstructure:"metrics"`
}

type (
	Listen struct {
		GatewayPort      string `mapstructure:"gateway_port"`
		GRPCPort         string `mapstructure:"grpc_port"`
		MigrationsPath   string `mapstructure:"migrations_path"`
		StocksServiceURL string `mapstructure:"stocks_service_url"`
		ServiceName      string `mapstructure:"service_name"`
		Env              string `mapstructure:"env"`
	}

	DbPostgres struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		DbName   string `mapstructure:"db_name"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Sslmode  string `mapstructure:"ssl_mode"`
	}

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	}

	Tracing struct {
		JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
	}

	Metrics struct {
		Port int64 `mapstructure:"port"`
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
