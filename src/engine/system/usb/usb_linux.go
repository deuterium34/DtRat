package usb

import "dtrat/errs"

type USBDevice struct {
	Name      string
	DeviceID  string
	VendorID  string
	ProductID string
	IsActive  bool
}

func GetUSBDevices() ([]USBDevice, error) {
	return nil, errs.ErrUnsupportedOs
}

func EnableDevice(deviceID string) error {
	return errs.ErrUnsupportedOs
}

// DisableDevice отключает устройство по его DeviceInstanceID
func DisableDevice(deviceID string) error {
	return errs.ErrUnsupportedOs
}

func EnableDevicePnPUtil(deviceID string) error {
	return errs.ErrUnsupportedOs
}

// DisableDevicePnPUtil резервный надежный метод отключения через консольную утилиту Windows
func DisableDevicePnPUtil(deviceID string) error {
	return errs.ErrUnsupportedOs
}
