package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
	DSN string `yaml:"dsn"`
}

type Server struct {
	Host string `yaml:"host"`
}

type Branding struct {
	Name string `yaml:"name"`
}

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Branding Branding `yaml:"branding"`
}

var loadedConfig *Config

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config *Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func init() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatal("error loading config file")
	}

	loadedConfig = config
}
