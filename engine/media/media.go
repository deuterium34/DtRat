package media

import (
	"dtrat/config"
	"fmt"
	"os"

	"github.com/hajimehoshi/go-mp3"
)

type Media struct {
	cfg config.Config
}

func NewMedia(cfg config.Config) (*Media, error) {
	m := Media{
		cfg: cfg,
	}
	return &m, nil
}

func (m *Media) Close() {

}

func (m *Media) Play(file *os.File) error {
	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		return fmt.Errorf("NewDecoder: %w", err)
	}

	pl, err := NewPlayer(false, m.cfg.Engine.SampleRate)
	if err != nil {
		return fmt.Errorf("NewPlayer: %w", err)
	}
	defer pl.Close()

	pl.LoadTrack(decodedMp3)
	err = pl.Play()
	if err != nil {
		return fmt.Errorf("Play: %w", err)
	}

	return nil
}
