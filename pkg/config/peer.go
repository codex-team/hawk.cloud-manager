package config

type Host struct {
	Name      string `yaml:"name" json:"name"`
	PublicKey string `yaml:"public_key" json:"public_key"`
}

type Group struct {
	Name string `yaml:"name" json:"name"`
	Hosts []string `yaml:"hosts" json:"hosts"`
}

type PeerConfig struct {
	Hosts  []Host `yaml:"hosts" json:"hosts"`
	Groups []Group `yaml:"groups" json:"groups"`
}
