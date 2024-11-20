package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	HTTP struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"http"`
	Postgres struct {
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
	} `yaml:"postgres"`
}

func MustLoad() *Config {
	file, err := os.Open("internal/config/config.yaml")
	if err != nil {
		panic("failed to open config file: " + err.Error())
	}
	defer file.Close()

	var cfg *Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		panic("failed to decode config file: " + err.Error())
	}

	return cfg
}
