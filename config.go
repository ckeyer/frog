package main

import (
	"os"
	"time"

	"github.com/ckeyer/commons/utils"
	yaml "gopkg.in/yaml.v2"
)

type Global struct {
	Period          utils.Duration
	DeleteEveryTime bool

	lastStart time.Time
}

type Registry struct {
	Name     string
	Username string
	Password string
	Server   string
}

type Task struct {
	Origin string
	Tags   []string
	Target string
}

type Config struct {
	Global     Global
	Registries []Registry
	Tasks      []Task
}

func OpenConfigFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
