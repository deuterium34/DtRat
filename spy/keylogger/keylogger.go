package keylogger

type Keylogger struct {
}

func NewKeylogger() (*Keylogger, error) {
	kl := Keylogger{}
	return &kl, nil
}
