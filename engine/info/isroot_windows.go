package info

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func isRoot() bool {
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
