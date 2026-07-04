package info

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func (i *Info) Report() string {
	var result strings.Builder

	// 1. Получаем текущего пользователя
	result.WriteString("Текущий пользователь: ")
	if u, err := user.Current(); err == nil {
		result.WriteString(u.Username + "\n")
	} else {
		result.WriteString(fmt.Sprintf("[Ошибка: %v]\n", err))
	}

	// 2. Получаем данные об OS из /etc/os-release
	result.WriteString("OS: ")
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		var prettyName string
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				prettyName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
				break
			}
		}
		if prettyName != "" {
			result.WriteString(prettyName + "\n")
		} else {
			result.WriteString("Linux (os-release найден, но имя не определено)\n")
		}
	} else {
		result.WriteString(fmt.Sprintf("[Ошибка чтения /etc/os-release: %v]\n", err))
	}

	// 3. Получаем версию ядра Linux через `uname -r`
	result.WriteString("Ядро: ")
	cmdKernel := exec.Command("uname", "-r")
	if out, err := cmdKernel.Output(); err == nil {
		result.WriteString(strings.TrimSpace(string(out)) + "\n")
	} else {
		result.WriteString("[Ошибка получения версии ядра]\n")
	}

	// 4. Архитектура через `uname -m`
	result.WriteString("Архитектура: ")
	cmdArch := exec.Command("uname", "-m")
	if out, err := cmdArch.Output(); err == nil {
		result.WriteString(strings.TrimSpace(string(out)) + "\n")
	} else {
		result.WriteString("[Ошибка получения архитектуры]\n")
	}

	return result.String()
}
