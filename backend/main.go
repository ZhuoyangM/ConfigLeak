package main

import (
	"errors"
	"fmt"
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

func main() {
	tests := []string{
		"https://stackoverflow.com/questions/5948659/when-should-i-use-a-trailing-slash-in-my-url",
		"http://chatgpt.com/c/6811c189-d6d4-8012-b2db-b8abcef1d053",
		"https://pkg.go.dev/net/url#URL",
		"http://cscsc.edu?q=123#abc",
	}

	paths, _ := LoadPathsFromFile("paths.yaml")
	fullUrls, err := BuildScanUrls(tests[3], paths)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Full URLs: ", fullUrls)
}
