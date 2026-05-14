package system

func (s *System) PressKey(key string) error {
	return ErrUnsupportedOs
}

func (s *System) Paste(text string) error {
	return ErrUnsupportedOs
}

func (s *System) PressHotKey(hotkey string) error {
	return ErrUnsupportedOs
}
