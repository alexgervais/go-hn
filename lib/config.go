package lib

import (
	"log"
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/go-homedir"
	"path"
)

type Config struct {
	EnableDebug bool
	BaseURL     string `yaml:"base-url"`
}

func LoadConfig() *Config {

	isDebug := flag.Bool("debug", false, "Display debug output")
	configFileLocation := flag.String("config", "", "Configuration file location")
	flag.Parse()

	if *configFileLocation == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Printf("Failed to retreive user home directory: %v", err)
		}
		*configFileLocation = path.Join(home, ".hn")
	}

	config := &Config{BaseURL: "https://hacker-news.firebaseio.com/v0", EnableDebug: *isDebug}

	data, err := ioutil.ReadFile(*configFileLocation)
	if err != nil {
		log.Printf("Failed to load configuration file: %v", err)
	}

	yaml.Unmarshal(data, config)

	return config
}
