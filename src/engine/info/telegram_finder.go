package info

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrNotFound = fmt.Errorf("папка Telegram не найдена")
)

// Поиск папки телеграмм по дефолтным путям
func (i *Info) findDefaultTelegramDir() (string, error) {
	var path string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("os.UserHomeDir: переменная окружения APPDATA не найдена")
		}
		path = filepath.Join(appData, "Telegram Desktop")

	case "darwin": // macOS
		path = filepath.Join(homeDir, "Library", "Application Support", "Telegram Desktop")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join(homeDir, "Library/Containers/ru.keepcoder.Telegram/Data/Library/Application Support/Telegram Desktop")
		}

	case "linux":
		path = filepath.Join(homeDir, ".local/share/TelegramDesktop")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join(homeDir, ".var/app/org.telegram.desktop/data/TelegramDesktop")
		}

	default:
		return "", fmt.Errorf("операционная система %s не поддерживается", runtime.GOOS)
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", ErrNotFound
	}
	if !info.IsDir() {
		return "", fmt.Errorf("найденный объект не является папкой: %s", path)
	}

	return path, nil
}

// Поиск папки телеграмм на дисках
func (i *Info) findGlobalTelegramDir() (string, error) {
	var foundPath string
	drives := i.sys.Drives()

	for _, root := range drives {
		if i.stopped.Load() {
			return "", fmt.Errorf("поиск остановлен")
		}
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				name := d.Name()
				if name == "Windows" || name == "$Recycle.Bin" || name == "System Volume Information" || name == "proc" || name == "dev" {
					return filepath.SkipDir
				}

				// Проверяем название папки
				if name == "Telegram Desktop" {
					foundPath = path
					return filepath.SkipAll
				}
			}
			return nil
		})

		if err != nil {
			return "", fmt.Errorf("Ошибка при обходе диска %s: %w", root, err)
		}
	}

	return foundPath, nil
}

func (i *Info) FindTelegramDir() (string, error) {
	if i.isTgPathFinding.Swap(true) {
		return "", fmt.Errorf("уже выполняется поиск папки Telegram")
	}
	defer i.isTgPathFinding.Store(false)

	path, err := i.findDefaultTelegramDir()
	if err != ErrNotFound {
		return path, nil
	}

	path, err = i.findGlobalTelegramDir()
	if err != nil {
		return "", ErrNotFound
	}

	return path, nil
}
