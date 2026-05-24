package info

import "os"

func isRoot() bool {
	if os.Geteuid() == 0 {
		return true
	} else {
		return false
	}
}
