package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB struct {
		Port     string `json:"port"`
		User     string `json:"user"`
		Database string `json:"database"`
		Password string `json:"password"`
	} `json:"db"`
	Port string `json:"port"`
}

func LoadConfiguration(filename string) (*Config, error) {
	config := new(Config)
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(configFile).Decode(config)
	return config, err
}
