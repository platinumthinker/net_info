package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddress string            `yaml:"listen_address"`
	Mikrotic      MikroticConfig    `yaml:"mikrotik"`
	StaticHosts   map[string]string `yaml:"static_hosts"`
}

type MikroticConfig struct {
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewConfig(cfgPath string) *Config {
	var cfg Config

	filename, _ := filepath.Abs(cfgPath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error read config %s: %v\n", cfgPath, err)
		return nil
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err == nil {
		return &cfg
	}

	fmt.Printf("Error parsing config %s: %v\n", cfgPath, err)

	return nil
}
