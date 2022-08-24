package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	CompareToRevision    string   `json:"compareToRevision"`
	TestRegex            string   `json:"testRegex"`
	IgnorePaths          []string `json:"ignorePaths"`
	ModuleFileExtensions []string `json:"moduleExtensions"`
	JestPath             string   `json:"jestPath"`
}

func LoadConfig(filename string) *Config {
	contents, err := os.ReadFile(filename)

	if errors.Is(err, os.ErrNotExist) {

		log.Println("- No config file found, creating a default file '" + filename + "'")

		// Doesn't exist, create the default config file
		jsonObject, _ := json.MarshalIndent(getDefaultConfig(), "", "    ")

		err = ioutil.WriteFile(filename, jsonObject, 0644)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Now it exists, load it
		return LoadConfig(filename)
	}

	// File exists, parse it
	var config Config
	json.Unmarshal(contents, &config)

	return &config
}

func getDefaultConfig() *Config {
	return &Config{
		JestPath:             `.\node_modules\.bin\jest`,
		TestRegex:            `(/_tests/.*|(\.|/)(test|spec))\.tsx?$`,
		IgnorePaths:          []string{"node_modules", ".idea", "coverage", ".git"},
		ModuleFileExtensions: []string{".js", ".jsx", ".ts", ".tsx"},
		CompareToRevision:    "origin/master",
	}
}
