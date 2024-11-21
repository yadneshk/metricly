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
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	configPath = "/etc/metricly/config.yaml"
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
}

func LoadConfig() (*Config, error) {
	conf, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("failed to open config file from default location: %s\n", configPath)
	}
	defer conf.Close()

	var cfg Config
	decoder := yaml.NewDecoder(conf)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %v", err)
	}
	return &cfg, nil
}
