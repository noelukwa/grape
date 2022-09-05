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

type Runner struct {
	Args    []string `json:"args"`
	Command string   `json:"command"`
}

type Namespace struct {
	Tag string

	Watch FWatcher `json:"watch"`

	Runner Runner `json:"run"`
}

type Config struct {
	Namespaces []Namespace
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
