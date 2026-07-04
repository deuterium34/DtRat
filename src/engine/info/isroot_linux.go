package info

import "os"

func (i *Info) IsRoot() bool {
	if os.Geteuid() == 0 {
		return true
	} else {
		return false
	}
}
