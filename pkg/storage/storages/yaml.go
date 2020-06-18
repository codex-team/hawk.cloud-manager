package storages

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlHost struct {
	Name      string `yaml:"name"`
	PublicKey string `yaml:"public_key"`
}

type YamlGroup struct {
	Name string `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

type YamlConfig struct {
	Hosts []YamlHost
	Groups []YamlGroup
}

type YamlStorage struct {
	Filename string
	Config YamlConfig
}

func NewYamlStorage(filename string) *YamlStorage {
	return &YamlStorage{Filename: filename}
}

func (s *YamlStorage) Load() error {
	yamlFile, err := ioutil.ReadFile(s.Filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &s.Config)
	if err != nil {
		return err
	}

	return nil
}

func (s *YamlStorage) Save() error {
	data, err := yaml.Marshal(s.Config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
