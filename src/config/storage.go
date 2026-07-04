package config

import (
	"fmt"
	"os"
)

type Storage interface {
	Get() (toml string, err error)
}

// ======================================

type fileStorage struct {
	path string
}

func (s *fileStorage) Get() (string, error) {
	file, err := os.ReadFile(s.path)
	if err != nil {
		return "", fmt.Errorf("os.ReadFile: %w", err)
	}

	return string(file), nil
}
