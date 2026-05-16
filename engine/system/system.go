package system

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

type System struct{}

func NewSystem() (*System, error) {
	s := System{}
	return &s, nil
}

func (s *System) Close() {

}

func (s *System) Exec(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("cmd.Run: %w", err)
	}

	return out.String(), nil
}

func (s *System) Screenshot() (string, error) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return "", fmt.Errorf("screenshot.NumActiveDisplays: Активные дисплеи не найдены")
	}

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", fmt.Errorf("screenshot.CaptureRect: %w", err)
	}

	filename := "screenshot.png"
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return "", fmt.Errorf("png.Encode: %w", err)
	}

	path, err := filepath.Abs(filename)
	if err != nil {
		return filename, nil
	}
	return path, nil
}

func (s *System) MonitorState(enabled bool) error {
	return setMonitorState(enabled)
}
