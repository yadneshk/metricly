/*
A sample yaml config file

server:

	address: 127.0.0.1
	port: 8080

prometheus:

	address: 127.0.0.1
	port: 9090

interval: 10
*/
package config

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	"gopkg.in/yaml.v3"
)

var (
	configPathDefault = "/etc/metricly/config.yaml"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
	Prometheus struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"prometheus"`
	CollectionInterval time.Duration `yaml:"interval"`
	Debug              bool          `yaml:"debug"`
}

func LoadConfig(configPath *string) (*Config, error) {
	if *configPath == "" {
		// if config file is not passed in cli, parse configs from configPathDefault
		configPath = &configPathDefault
	}

	// attempt to load config
	conf, err := os.Open(*configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file from default location: %s", *configPath)
	}
	defer conf.Close()

	var cfg Config
	decoder := yaml.NewDecoder(conf)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %v", err)
	}
	slog.Info(fmt.Sprintf("Config loaded successfully from %s", *configPath))

	// checking if any variable was overrided through environment variables
	if env := os.Getenv("SERVER_ADDRESS"); env != "" {
		cfg.Server.Address = env
	}
	if env := os.Getenv("SERVER_PORT"); env != "" {
		cfg.Server.Port = env
	}
	if env := os.Getenv("PROMETHEUS_ADDRESS"); env != "" {
		cfg.Prometheus.Address = env
	}
	if env := os.Getenv("PROMETHEUS_PORT"); env != "" {
		cfg.Prometheus.Port = env
	}
	if env := os.Getenv("COLLECTION_INTERVAL"); env != "" {
		if interval, err := time.ParseDuration(env); err == nil {
			cfg.CollectionInterval = interval
		} else {
			return nil, fmt.Errorf("invalid COLLECTION_INTERVAL value: %v", err)
		}
	}
	if env := os.Getenv("DEBUG"); env != "" {
		if debug, err := parseBool(env); err == nil {
			cfg.Debug = debug
		} else {
			return nil, fmt.Errorf("invalid DEBUG value: %v", err)
		}
	}

	return &cfg, nil
}

// parseBool parses a string into a boolean value.
func parseBool(value string) (bool, error) {
	switch value {
	case "1", "true", "TRUE", "True":
		return true, nil
	case "0", "false", "FALSE", "False":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}
