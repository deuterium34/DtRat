package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	/*
		Доступные хранилища:
		file
	*/
	UseStorage = "file"

	Version = "0.0.0"
)

var (
	ErrUndefainedStorage = errors.New("Неизвестное хранилище")
)

// ======================================

type Config struct {
	General GenaralConfig `toml:"General"`
	Bot     BotConfig     `toml:"Bot"`
	Engine  EngineConfig  `toml:"Engine"`
}

type BotConfig struct {
	Token  string `toml:"Token"`
	UserID int    `toml:"User_id"`
}

type EngineConfig struct {
	SampleRate int `toml:"Sample_rate"`
}

type GenaralConfig struct {
	HostName string `toml:"Host_name"`
}

// ======================================

func NewConfig() (Config, error) {
	var storage Storage
	switch UseStorage {
	case "file":
		storage = &fileStorage{path: "config.toml"}
	default:
		return Config{}, ErrUndefainedStorage
	}

	cfgRaw, err := storage.Get()
	if err != nil {
		return Config{}, fmt.Errorf("storage.Get: %w", err)
	}

	var cfg Config
	_, err = toml.Decode(cfgRaw, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("Decode: %w", err)
	}

	return cfg, nil
}
