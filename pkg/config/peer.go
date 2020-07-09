package config

type Host struct {
	Name      string `yaml:"name"`
	PublicKey string `yaml:"public_key"`
}

type Group struct {
	Name string `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

type PeerConfig struct {
	Hosts  []Host `yaml:"hosts"`
	Groups []Group `yaml:"groups"`
}
