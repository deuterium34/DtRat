package config

const (
	hardcodedConfig = ``
)

type hardcodedStorage struct{}

func (s *hardcodedStorage) Get() (string, error) {
	return hardcodedConfig, nil
}
