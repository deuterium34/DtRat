package sniffer

type Sniffer struct {
}

func NewSniffer() (*Sniffer, error) {
	sn := Sniffer{}
	return &sn, nil
}
