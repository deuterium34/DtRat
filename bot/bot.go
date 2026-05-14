package bot

import (
	"errors"
	"strings"
)

var (
	ErrWaitAborted = errors.New("Ожидание сообщения было прервано")
)

type Bot interface {
	// Методы для отправки сообщений
	Send(s string, args ...any) error
	SendFile(file string) error

	// Методы для получения сообщений
	Wait() (message string, err error)
	WaitFile() (path string, err error)

	// Методы для упровления
	Start() error
	Close() error
}

// Парсинг строки "/cmd arg1 arg2 arg3 ..."
func ParseCommand(input string) (command, args string) {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "/")

	parts := strings.SplitN(input, " ", 2)

	if len(parts) == 0 {
		return "", ""
	}

	command = parts[0]

	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}

	return command, args
}
