package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type PathConfig struct {
	Paths []string `yaml:"paths"`
}

func LoadPathsFromFile(filepath string) ([]string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("error reading yaml file")
	}

	var config PathConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New("error unmarshalling yaml file")
	}

	return config.Paths, nil
}

func ParseBaseUrl(url string) (string, error) {
	// TODO
	return "", nil
}

func BuildScanUrls(base string, paths []string) ([]string, error) {
	// TODO
	return nil, nil
}

func main() {
	paths, err := LoadPathsFromFile("paths.yaml")
	if err != nil {
		fmt.Println("Error loading paths:", err)
		return
	}
	fmt.Println("Loaded paths:", paths)
	fmt.Println("Count: ", len(paths))
}
