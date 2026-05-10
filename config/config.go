package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	Version = "0.0.0"
)

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

func NewConfig(configPath string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("DecodeFile: %w", err)
	}

	return cfg, nil
}
