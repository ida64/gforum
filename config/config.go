package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
	DSN             string `yaml:"dsn"`
	UseSqlite       bool   `yaml:"useSqlite"`
	SqlitePath      string `yaml:"sqlitePath"`
	MigrateToSqlite bool   `yaml:"migrateToSqlite"`
}

type Server struct {
	Host string `yaml:"host"`
}

type Branding struct {
	Name string `yaml:"name"`
}

type Category struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Config struct {
	Database   Database   `yaml:"database"`
	Server     Server     `yaml:"server"`
	Branding   Branding   `yaml:"branding"`
	Categories []Category `yaml:"categories"`
}

var LoadedConfig *Config

/*
* loadConfig parses the config from the specified path and returns it.
* It returns an error if the config could not be loaded.
 */
func LoadConfig(path string) (*Config, error) {
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

/*
* saveConfig saves the supplied config to the specified path.
* It returns an error if the config could not be saved.
 */
func SaveConfig(path string, config *Config) error {
	var data, err = yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("error loading config file")
	}

	LoadedConfig = config
}
