package yaml

import (
	"io/ioutil"

	"github.com/codex-team/hawk.cloud-manager/internal/storage"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"gopkg.in/yaml.v2"
)

// YamlStorage is an implementation of Storage interface that stores peer
// configuration in a yaml file
type YamlStorage struct {
	// Name of file to store Peer Config in
	Filename string
	// Current Peer configuration
	config config.PeerConfig
}

// NewYamlStorage creates YamlStorage instanse and returns it as a Storage
// interface implementation
func NewYamlStorage(filename string) storage.Storage {
	return &YamlStorage{Filename: filename}
}

// Load reads Peer Config from file and saves it to `config` field
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

// Get returns stored Peer Config
func (s *YamlStorage) Get() *config.PeerConfig {
	return &s.config
}

// Save writes Peer Config from `config` field to file
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
