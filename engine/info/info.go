package info

type Info struct {
}

func NewInfo() (*Info, error) {
	i := Info{}
	return &i, nil
}

func (i *Info) Close() {

}
