package media

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/ebitengine/oto/v3"
)

type Player struct {
	otoCtx *oto.Context
	player *oto.Player
	mu     sync.Mutex
}

var ErrNotLoaded error = errors.New("Track not loaded")

func NewPlayer(mono bool, sampleRate int) (*Player, error) {
	var ch int
	if mono {
		ch = 1
	} else {
		ch = 2
	}

	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: ch,
		Format:       oto.FormatSignedInt16LE,
	}

	octx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, fmt.Errorf("oto.NewContext: %w", err)
	}
	<-readyChan

	return &Player{
		otoCtx: octx,
		player: nil,
	}, nil
}

func (p *Player) LoadTrack(r io.Reader) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.player = p.otoCtx.NewPlayer(r)
}

func (p *Player) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.player == nil {
		return ErrNotLoaded
	}

	p.player.Play()

	return nil
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.player == nil {
		return
	}

	p.player.Pause()
}

func (p *Player) PlayPause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.player == nil {
		return
	}

	if p.player.IsPlaying() {
		p.player.Pause()
	} else {
		p.player.Play()
	}
}

func (p *Player) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.player == nil {
		return
	}

	p.player.Pause()
	p.player = nil
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.player == nil {
		return false
	}
	return p.player.IsPlaying()
}
