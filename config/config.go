package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Global struct {
	My   MySQL  `yaml:"mysql"`
	Serv Server `yaml:"server"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	Glb  Global
	My   MySQL
	Serv Server
)

// Load config from file
func Load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Println("Error: Config file failed to read file -", err)
		return err
	}

	err = yaml.Unmarshal(data, &Glb)
	if err != nil {
		log.Println("Error: Config file failed to Unmarshal -", err)
		return err
	}

	My = Glb.My
	Serv = Glb.Serv

	return nil
}
