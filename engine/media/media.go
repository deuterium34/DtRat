package media

import (
	"dtrat/config"
	"dtrat/errs"
	"fmt"
	"os"
	"os/exec"
	"runtime"

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

func (m *Media) OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errs.ErrUnsupportedOs
	}

	if err != nil {
		return err
	}
	return nil
}
