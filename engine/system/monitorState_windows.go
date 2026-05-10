package system

import (
	"syscall"
)

const (
	WM_SYSCOMMAND   = 0x0112
	SC_MONITORPOWER = 0xF170
	HWND_BROADCAST  = 0xFFFF
)

func setMonitorState(enabled bool) error {
	var state uintptr
	if enabled {
		state = ^uintptr(0)
	} else {
		state = 2
	}

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessage := user32.NewProc("SendMessageW")

	_, _, err := sendMessage.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SYSCOMMAND),
		uintptr(SC_MONITORPOWER),
		state,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return err
	}
	return nil
}
