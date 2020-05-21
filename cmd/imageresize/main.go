package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/assanoff/image-resizer/internal/apiserver"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.yaml", "path to config file")
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

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
