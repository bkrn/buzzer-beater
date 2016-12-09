//Package config interprets and buils config file
package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type databaseConfig struct {
	Name    string
	Buckets []string
}

type application struct {
	Name   string
	Secret string
	Salt   string
}

type owner struct {
	Name  string
	Email string
}

type tomlConfig struct {
	Title       string
	Owner       owner
	Application application
	Database    databaseConfig
}

func getConfig() tomlConfig {
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config
}

//Config data for global sharing, yay!
var Config = getConfig()
