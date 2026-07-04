package system

import (
	"fmt"
	"os/exec"
)

func (s *System) SetVolume(percent int) error {
	if percent < 0 || percent > 100 {
		return fmt.Errorf("volume percent must be between 0 and 100")
	}

	// Форматируем строку в вид "50%"
	volumeStr := fmt.Sprintf("%d%%", percent)

	// Пробуем использовать pactl (PulseAudio / PipeWire) — стандарт для современных Linux-десктопов
	cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", volumeStr)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Если pactl не найден или выдал ошибку, пробуем amixer (ALSA)
	cmd = exec.Command("amixer", "-D", "pulse", "sset", "Master", volumeStr)
	if err := cmd.Run(); err != nil {
		// Если и amixer не сработал, пробуем стандартный Master канал ALSA
		cmd = exec.Command("amixer", "sset", "Master", volumeStr)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set volume via pactl and amixer: %w", err)
		}
	}

	return nil
}
