package hider

type Hider struct {
}

func NewHider() (*Hider, error) {
	hdr := Hider{}
	return &hdr, nil
}
