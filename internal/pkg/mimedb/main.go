package mimedb

import (
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

const configPath = "./configs/mime.yml"

var Config MimeConfig
var once sync.Once

func TypeByExtension(ext string) (mime string, found bool) {
	once.Do(loadTypes)

	m, f := Config.Extensions[ext]
	if !f {
		return Config.Default, false
	}
	return m, true
}

func loadTypes() {
	var c MimeConfig
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		panic(err)
	}

	Config = c
}

type MimeConfig struct {
	Default    string            `yaml:"default"`
	Extensions map[string]string `yaml:"extensions"`
}
