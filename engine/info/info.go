package info

import (
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Info struct {
}

func NewInfo() (*Info, error) {
	i := Info{}
	return &i, nil
}

func (i *Info) IsRoot() bool {
	switch runtime.GOOS {
	case "windows":
		return isRootWin()
	case "linux":
		return isRootLinux()
	default:
		return false
	}
}

func isRootLinux() bool {
	if os.Geteuid() == 0 {
		return true
	} else {
		return false
	}
}

func isRootWin() bool {
	var token windows.Token

	//err := windows.OpenCurrentProcessToken(windows.GetCurrentProcess(), windows.TOKEN_QUERY, &token)
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	var elevation uint32
	var returnLength uint32
	err = windows.GetTokenInformation(
		token,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&returnLength,
	)

	if err != nil {
		return false
	}

	return elevation != 0
}
