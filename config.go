package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Port    string
	Limit int
}

func ReadConfig(path string) (*Configuration, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		return nil, err
	}

	return &configuration, nil
}