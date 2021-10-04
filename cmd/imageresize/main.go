package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/assanoff/image-resizer/internal/apiserver"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config-path", "config.yaml", "path to config file")
	port := os.Getenv("PORT")
	debugLevel := os.Getenv("DEBUG_LEVEL")

	flag.Parse()

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Errorf("could not read config file", err)
	}

	config := apiserver.NewConfig()
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		fmt.Errorf("could not unmarshal config", err)
	}
	config.Port = ":" + port
	config.LogLevel = debugLevel

	log.Printf("Starting on %s", config.Port)
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
