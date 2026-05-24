package info

import (
	"dtrat/engine/system"
	"sync/atomic"
)

type Info struct {
	sys system.System

	stopped         atomic.Bool
	isTgPathFinding atomic.Bool
}

func NewInfo(system *system.System) (*Info, error) {
	i := Info{
		sys: *system,
	}
	return &i, nil
}

func (i *Info) Stop() {
	i.stopped.Store(true)
}
