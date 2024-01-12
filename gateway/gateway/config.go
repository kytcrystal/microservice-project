package gateway

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Routes []Route
}

type Route struct {
	Prefix string
	Host   string
}

func ReadConfigFromFile(filename string) (*Configuration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var configuration Configuration
	err = yaml.NewDecoder(f).Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}
