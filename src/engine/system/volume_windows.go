package system

import (
	"fmt"
	"os/exec"
	"syscall"
)

func (s *System) SetVolume(percent int) error {
	if percent < 0 || percent > 100 {
		return fmt.Errorf("volume percent must be between 0 and 100")
	}

	// Но так как чистый .NET класс [audio] доступен не во всех старых версиях Windows по умолчанию,
	// самый надежный и гарантированный способ управления через стандартный SoundDevice в PS:
	script := fmt.Sprintf(
		`(New-Object -ComObject ScriptControl).Language = 'VBScript'; `+
			`$wsh = New-Object -ComObject WScript.Shell; `+
			`[void]$wsh.AppActivate('Volume'); `+
			`$psvol = (New-Object -ComObject MMDeviceEnumerator).GetDefaultAudioEndpoint(0,0).AudioEndpointVolume; `+
			`$psvol.MasterVolumeLevelScalar = %f`, float64(percent)/100.0,
	)

	cmd := exec.Command("powershell", "-Command", script)

	// Скрываем окно консоли PowerShell при вызове из GUI-приложения
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set volume via PowerShell: %w", err)
	}

	return nil
}
