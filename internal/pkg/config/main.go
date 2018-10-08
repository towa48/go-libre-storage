package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const configPath = "./configs/default.yml"

var config *Config

type Config struct {
	UsersDb string `yaml:"usersDb"`
	FilesDb string `yaml:"filesDb"`
	Storage string `yaml:"storage"`
}

func Get() Config {
	if config != nil {
		return *config
	}

	var c Config
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		panic(err)
	}

	//config = new(Config)
	config = &c
	return *config
}
