package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	/*
		Доступные хранилища:
		file, hardcoded
	*/
	UseStorage = "file"

	Version = "0.0.0"
)

var (
	ErrUndefainedStorage = errors.New("Неизвестное хранилище")
)

// ======================================
// Структуры для хранения конфигурации, соответствующие структуре TOML-файла

type Config struct {
	General   GeneralConfig   `toml:"General"`
	Transport TransportConfig `toml:"Transport"`
	Engine    EngineConfig    `toml:"Engine"`
}

type EngineConfig struct {
	SampleRate int `toml:"Sample_rate"` // Частота дискретизации для аудио
}

type GeneralConfig struct {
	AgentName    string `toml:"Agent_name"`
	UseTransport string `toml:"Use_transport"`
}

type TransportConfig struct {
	Telegram *TelegramConfig `toml:"Telegram"`
	Arcanum  *ArcanumConfig  `toml:"Arcanum"`
}

type TelegramConfig struct {
	Token  string `toml:"Token"`
	UserID int    `toml:"User_id"`
}

type ArcanumConfig struct {
	Addr   string `toml:"Addr"`
	Secret string `toml:"Secret"`
}

// ======================================

func NewConfig() (Config, error) {
	var storage Storage
	switch UseStorage {
	case "file":
		storage = &fileStorage{path: "config.toml"}
	case "hardcoded":
		storage = &hardcodedStorage{}
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
