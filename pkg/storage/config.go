package storage

type Host struct {
	Name string
	PublicKey string
}

type Group struct {
	Name string
	Hosts []string
}

type Config struct {
	Hosts []Host
	Groups []Group
}