package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Listen_port string   `yaml:"listen_port"`
	Log_file    string   `yaml:"log_file"`
	Tg_api      string   `yaml:"tg_api"`
	Tg_bot_id   string   `yaml:"tg_bot_id"`
	Track_dirs  []string `yaml:"track_dirs"`
	Track_files []string `yaml:"track_files"`
}

func (c *Conf) Read() (*Conf, error) {

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return c, err
	}

	return c, nil
}
