package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type FWatcher struct {
	Exclude []string `json:"exclude"`

	Include []string `json:"include"`
}

type Namespace struct {
	Tag string

	Watch FWatcher `json:"watch"`

	Run string `json:"run"`
}

type Config struct {
	Namespaces []Namespace
}

func (conf *Config) getNameSpace(tag string) *Namespace {
	for _, ns := range conf.Namespaces {

		if ns.Tag == tag {
			return &ns
		}
	}
	return nil
}

func ConfigFromJson(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	if filepath.Ext(abs) != ".json" {
		return nil, errors.New("only json config files are supported")
	}

	jsonDump, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = json.Unmarshal(jsonDump, &config)
	if err != nil {
		return nil, err
	}

	var namespaces []Namespace
	for key, value := range config {
		var ns Namespace
		ns.Tag = key
		vJson, _ := json.Marshal(value)
		err = json.Unmarshal(vJson, &ns)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}

	return &Config{Namespaces: namespaces}, nil
}

func ConfigFromFlags() (*Config, error) {
	return nil, nil
}

func createConfigFile() error {

	var defaultConfig = []byte(`{
	"dev": {
		"watch": {
			"exclude": [
			  "vendor"
			],
		"include": [
				"*.go"
			 ]
		  },
		  "run":  "go run ."
		}
	  }`)

	//create a json file and write the contents of defaultConfig to it
	if err := os.WriteFile("grape.json", defaultConfig, 0644); err != nil {
		return err
	}
	return nil
}
