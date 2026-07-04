package system

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rpdg/winput"
)

var (
	ErrUndefinedRune    = errors.New("Неизвестный символ")
	ErrUndefinedKeyName = errors.New("Неизвестная клавиша")
)

func keyFromName(keyName string) (winput.Key, error) {
	keyName = strings.ToLower(strings.TrimSpace(keyName))

	switch keyName {
	// Управление и ввод
	case "enter", "return":
		return winput.KeyEnter, nil
	case "esc", "escape":
		return winput.KeyEsc, nil
	case "space", "spc":
		return winput.KeySpace, nil
	case "tab":
		return winput.KeyTab, nil
	case "backspace", "bs":
		return winput.KeyBkSp, nil
	case "ins", "insert":
		return winput.KeyInsert, nil
	case "del", "delete":
		return winput.KeyDelete, nil

	// Модификаторы
	case "shift":
		return winput.KeyShift, nil
	case "ctrl", "control":
		return winput.KeyCtrl, nil
	case "alt", "menu":
		return winput.KeyAlt, nil
	//case "lwin", "win", "cmd": return winput.Key, nil
	//case "rwin":            return winput.KeyRWin, nil
	case "capslock":
		return winput.KeyCaps, nil

	// Навигация
	case "up":
		return winput.KeyArrowUp, nil
	case "down":
		return winput.KeyArrowDown, nil
	case "left":
		return winput.KeyLeft, nil
	case "right":
		return winput.KeyRight, nil
	case "home":
		return winput.KeyHome, nil
	case "end":
		return winput.KeyEnd, nil
	case "pgup", "pageup":
		return winput.KeyPageUp, nil
	case "pgdn", "pagedown":
		return winput.KeyPageDown, nil

	// Функциональные клавиши
	case "f1":
		return winput.KeyF1, nil
	case "f2":
		return winput.KeyF2, nil
	case "f3":
		return winput.KeyF3, nil
	case "f4":
		return winput.KeyF4, nil
	case "f5":
		return winput.KeyF5, nil
	case "f6":
		return winput.KeyF6, nil
	case "f7":
		return winput.KeyF7, nil
	case "f8":
		return winput.KeyF8, nil
	case "f9":
		return winput.KeyF9, nil
	case "f10":
		return winput.KeyF10, nil
	case "f11":
		return winput.KeyF11, nil
	case "f12":
		return winput.KeyF12, nil

	default:
		return 0, ErrUndefinedKeyName
	}
}

func keyFromString(key string) (winput.Key, error) {
	var k winput.Key
	char := []rune(key)[0]

	if len(key) == 1 {
		var ok bool
		k, ok = winput.KeyFromRune(char)
		if !ok {
			return 0, ErrUndefinedRune
		}
	} else {
		var err error
		k, err = keyFromName(key)
		if err != nil {
			return 0, fmt.Errorf("keyFromName: %w", err)
		}
	}

	return k, nil
}

func (s *System) PressKey(key string) error {
	k, err := keyFromString(key)
	if err != nil {
		return fmt.Errorf("keyFromString: %w", err)
	}

	err = winput.Press(k)
	if err != nil {
		return fmt.Errorf("winput.Press: %w", err)
	}
	return nil
}

func (s *System) Paste(text string) error {
	err := winput.Type(text)
	if err != nil {
		return fmt.Errorf("winput.Type: %w", err)
	}
	return nil
}

// нажатие хоткеев key+key
func (s *System) PressHotKey(hotkey string) error {
	hotkey = strings.TrimSpace(hotkey)
	hotkeys := strings.Split(hotkey, "+")

	keys := make([]winput.Key, len(hotkeys))

	for i, key := range hotkeys {
		k, err := keyFromString(key)
		if err != nil {
			return fmt.Errorf("keyFromString (%s): %w", key, err)
		}
		keys[i] = k
	}

	err := winput.PressHotkey(keys...)
	if err != nil {
		return fmt.Errorf("winput.PressHotkey: %w", err)
	}
	return nil
}
