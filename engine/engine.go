package engine

import (
	"dtrat/config"
	"dtrat/engine/info"
	"dtrat/engine/media"
	"dtrat/engine/system"

	"fmt"
)

type Engine struct {
	System *system.System
	Info   *info.Info
	Media  *media.Media
}

func NewEngine(cfg config.Config) (*Engine, error) {
	sys, err := system.NewSystem()
	if err != nil {
		return nil, fmt.Errorf("NewSystem: %w", err)
	}

	inf, err := info.NewInfo()
	if err != nil {
		return nil, fmt.Errorf("NewInfo: %w", err)
	}

	med, err := media.NewMedia(cfg)
	if err != nil {
		return nil, fmt.Errorf("NewMedia: %w", err)
	}

	eng := Engine{
		System: sys,
		Info:   inf,
		Media:  med,
	}

	return &eng, nil
}

func (e *Engine) Close() {
	e.Info.Close()
	e.Media.Close()
	e.System.Close()
}
