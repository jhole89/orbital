package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Database *DatabaseConfig `yaml:"database"`
	Service  *ServiceConfig  `yaml:"service"`
	Lakes    []*LakeConfig   `yaml:"lakes"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Endpoint string `yaml:"endpoint"`
}

type ServiceConfig struct {
	Port int64 `yaml:"port"`
}

type LakeConfig struct {
	Provider string `yaml:"provider"`
	Store    string `yaml:"store"`
	Address  string `yaml:"address"`
}

func (c *Config) getConf() *Config {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Unable to read config.yaml: %s\n", err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unable to expand config.yaml: %s\n", err.Error())
	}

	dbAddress, ok := os.LookupEnv("ORBITAL_DB_ADDRESS")
	if ok {
		if dbAddress != "" {
			c.Database.Endpoint = dbAddress
		}
	}

	return c
}
