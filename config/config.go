package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type config struct {
	MailHost string `json"mailHost"`
	MailPort int    `json"mailPort"`
	MailUser string `json"mailUser"`
	MailPass string `json"mailPass"`
	DbHost   string `json"dbHost"`
	DbName   string `json"dbName"`
}

var instance *config
var once sync.Once

func loadConfiguration(file string) *config {
	var conf *config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&conf)
	return conf
}

func GetInstance() *config {
	once.Do(func() {
		instance = loadConfiguration("./config.json")
	})
	return instance
}
