package info

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
)

func (i *Info) Report() string {
	var result strings.Builder

	// 1. Получаем текущего пользователя (пакет os/user в Go изначально работает в UTF-8)
	result.WriteString("Текущий пользователь: ")
	if u, err := user.Current(); err == nil {
		result.WriteString(u.Username + "\n")
	} else {
		result.WriteString(fmt.Sprintf("[Ошибка: %v]\n", err))
	}

	// Команда для PowerShell, которая принудительно устанавливает UTF-8 для вывода
	utf8Prefix := "[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; "

	// 2. Получаем данные об OS
	result.WriteString("OS: ")

	psScript := utf8Prefix + "(Get-CimInstance Win32_OperatingSystem).Caption + ' ' + (Get-CimInstance Win32_OperatingSystem).DisplayVersion"
	cmd := exec.Command("powershell", "-Command", psScript)

	// Скрываем окно консоли
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		osInfo := strings.TrimSpace(out.String())
		if osInfo != "" {
			result.WriteString(osInfo + "\n")
		} else {
			result.WriteString("Windows (Не удалось определить точную версию)\n")
		}
	} else {
		result.WriteString(fmt.Sprintf("[Ошибка получения OS: %v]\n", err))
	}

	// 3. Архитектура
	result.WriteString("Архитектура: ")
	cmdArch := exec.Command("powershell", "-Command", utf8Prefix+"(Get-CimInstance Win32_OperatingSystem).OSArchitecture")
	cmdArch.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out.Reset()
	cmdArch.Stdout = &out
	if err := cmdArch.Run(); err == nil {
		result.WriteString(strings.TrimSpace(out.String()) + "\n")
	} else {
		result.WriteString("[Ошибка получения архитектуры]\n")
	}

	return result.String()
}
