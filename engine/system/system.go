package system

import (
	"fmt"
	"image/png"
	"os"
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
