package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Jwt      Jwt      `yaml:"jwt"`
}

type Server struct {
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

type Database struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schemas  string `yaml:"schemas"`
}

type Jwt struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expireTime"`
}

var Setting = &Config{}

func init() {
	yamlFile, err := ioutil.ReadFile("./conf/web.yml")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.UnmarshalStrict(yamlFile, Setting)
	if err != nil {
		log.Fatal(err.Error())
	}
}
