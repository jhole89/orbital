package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Service ServiceConfig `yaml:"service"`
	Lakes []LakeConfig `yaml:"lakes"`
}

type DatabaseConfig struct {
	Endpoint string `yaml:"endpoint"`
}

type ServiceConfig struct {
	Port int64 `yaml:"port"`
}

type LakeConfig struct {
	Provider string `yaml:"provider"`
	Store string `yaml:store`
	Address string `yaml:address`
}

func (c *Config) getConf() *Config {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}