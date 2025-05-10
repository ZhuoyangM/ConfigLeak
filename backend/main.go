package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
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
	// TODO
	return nil, nil
}

func main() {
	tests := []string{
		"https://stackoverflow.com/questions/5948659/when-should-i-use-a-trailing-slash-in-my-url",
		"http://chatgpt.com/c/6811c189-d6d4-8012-b2db-b8abcef1d053",
		"https://pkg.go.dev/net/url#URL",
		"http://cscsc.edu?q=123#abc",
	}

	for _, test := range tests {
		res, _ := ParseBaseUrl(test)
		fmt.Println(res)
	}
}
