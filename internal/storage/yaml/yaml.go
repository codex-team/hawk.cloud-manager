package yaml

import (
	"io/ioutil"

	"github.com/codex-team/hawk.cloud-manager/internal/storage"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"gopkg.in/yaml.v2"
)

type YamlStorage struct {
	Filename string
	config   config.PeerConfig
}

func NewYamlStorage(filename string) storage.Storage {
	return &YamlStorage{Filename: filename}
}

func (s *YamlStorage) Load() error {
	yamlFile, err := ioutil.ReadFile(s.Filename)
	if err != nil {
		return err
	}

	s.config = config.PeerConfig{}

	err = yaml.Unmarshal(yamlFile, &s.config)
	if err != nil {
		return err
	}

	return nil
}

func (s *YamlStorage) Get() *config.PeerConfig {
	return &s.config
}

func (s *YamlStorage) Save() error {
	data, err := yaml.Marshal(s.config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
