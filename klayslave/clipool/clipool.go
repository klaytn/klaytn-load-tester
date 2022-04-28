package clipool

import "sync"

type ClientCreatorFunc func() interface{}

type ClientPool struct {
	max  int
	init int
	cnt  int

	freeList []interface{}

	lock      sync.Mutex
	allocFunc ClientCreatorFunc
}

func (p *ClientPool) Init(init int, max int, allocFunc ClientCreatorFunc) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.init = init
	p.max = max
	p.allocFunc = allocFunc
	for i := 0; i < p.init; i++ {
		cli := p.allocFunc()
		p.freeList = append(p.freeList, cli)
		p.cnt++
	}
}

func (p *ClientPool) Alloc() interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()

	var ret interface{}
	if len(p.freeList) == 0 {
		cli := p.allocFunc()
		ret = cli
		p.cnt++
	} else {
		ret = p.freeList[0]
		p.freeList = p.freeList[1:]
	}
	return ret
}

func (p *ClientPool) Free(v interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.freeList = append(p.freeList, v)
}
