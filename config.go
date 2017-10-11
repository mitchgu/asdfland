package main

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}

	v.SetDefault("server_addr", "localhost:9090")
	v.SetDefault("db_kind", "redis")
	v.SetDefault("redis_addr", "localhost:6379")
	v.SetDefault("redis_db", 0)
	v.SetDefault("redis_pass", "")
	v.SetDefault("frontend_dir", "frontend")

	return v
}

var c = LoadConfig()
