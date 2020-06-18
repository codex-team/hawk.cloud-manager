package storage

type IStorage interface {
	Load() error
	Save() error
}

type Storage struct {
	Config Config
}