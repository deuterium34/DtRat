package rat

var helpText = `
Полная документация на https://github.com/deuterium34/DtRat

[W] - Доступно только в Windows
Справка по доступным командам:
/start - Запускает первичную инициализацию (требуется один раз)
/help - Показывает эту справку
/kill - Завершает работу 
/screenshot - Скриншот основного экрана
/monitor [on|off] - [W] Управляет состоянием монитора
/keyboard [paste|press|hotkey] STRING - [W] Симулирует ввод с клавиатуры
	paste - Вставляет текст
	press - Одиночное нажатие клавиши
	hotkey - Нажимает сочитание клавишь (+ в качестве разделителя)
/browser STRING - Открывает URL в браузере
/findtg - Выполняет поиск папки Telegram
/usb [list|enable|disable] - [W] Управление USB устройствами
	list - Показывает список USB устройств
	enable [ID] - Включает USB устройство по ID из списка
	disable [ID] - Отключает USB устройство по ID из списка
`
