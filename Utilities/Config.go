package Utilities

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	ApiTitle    string `json:"apiTitle"`
	ApiVersion  string `json:"apiVersion"`
	RedBaseUrl  string `json:"redBaseUrl"`
	BlueBaseUrl string `json:"blueBaseUrl"`
}

var GlobalConfig Config

func GetConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
