package post

type ChanPair struct {
	In  chan interface{}
	Out chan interface{}
}

func (ch ChanPair) Init(size int) ChanPair {
	ch.In = make(chan interface{}, size)
	ch.Out = make(chan interface{}, size)
	return ch
}
