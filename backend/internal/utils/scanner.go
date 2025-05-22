package utils

import (
	"errors"
	"net/url"
	"os"
	"path"
	"strings"

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

func ParseBaseUrl(input string) (string, error) {
	if !strings.HasPrefix(input, "http") && !strings.HasPrefix(input, "https") {
		return "", errors.New("url must start with http or https")
	}
	u, err := url.Parse(input)
	if err != nil {
		return "", errors.New("error parsing url")
	}
	u.Path = "/"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func BuildScanUrls(base string, paths []string) ([]string, error) {
	parsedBase, err := ParseBaseUrl(base)
	if err != nil {
		return nil, err
	}

	baseUrl, _ := url.Parse(parsedBase)
	var fullUrls []string
	for _, p := range paths {
		fullUrl := *baseUrl
		fullUrl.Path = path.Join(baseUrl.Path, p)
		fullUrls = append(fullUrls, fullUrl.String())
	}
	return fullUrls, nil
}
