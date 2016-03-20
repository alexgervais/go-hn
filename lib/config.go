package lib

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/go-homedir"
	"path"
)

type Config struct {
	BaseURL string `yaml:"base-url"`
}

func LoadConfig(configFileLocation string) (config *Config) {

	config = &Config{BaseURL: "https://hacker-news.firebaseio.com/v0"}

	if configFileLocation == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Printf("Failed to retreive user home directory: %v", err)
			return
		}
		configFileLocation = path.Join(home, ".hn")
	}

	data, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		log.Printf("Failed to load configuration file: %v", err)
		return
	}

	yaml.Unmarshal(data, config)
	return
}
