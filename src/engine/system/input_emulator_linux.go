package system

import "dtrat/errs"

func (s *System) PressKey(key string) error {
	return errs.ErrUnsupportedOs
}

func (s *System) Paste(text string) error {
	return errs.ErrUnsupportedOs
}

func (s *System) PressHotKey(hotkey string) error {
	return errs.ErrUnsupportedOs
}
