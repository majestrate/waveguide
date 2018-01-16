package streaming

import "sync"

type Context struct {
	mtx     sync.Mutex
	streams map[string]*StreamInfo
}

func NewContext() *Context {
	return &Context{
		streams: make(map[string]*StreamInfo),
	}
}

func (ctx *Context) Online() (keys []string) {
	ctx.mtx.Lock()
	for k := range ctx.streams {
		keys = append(keys, k)
	}
	ctx.mtx.Unlock()
	return
}

func (ctx *Context) Ensure(k string) (i *StreamInfo) {
	if len(k) > 0 {
		ctx.mtx.Lock()
		var ok bool
		i, ok = ctx.streams[k]
		if !ok {
			i = new(StreamInfo)
			ctx.streams[k] = i
		}
		ctx.mtx.Unlock()
	}
	return
}

func (ctx *Context) Find(k string) (i *StreamInfo) {
	if len(k) > 0 {
		ctx.mtx.Lock()
		i, _ = ctx.streams[k]
		ctx.mtx.Unlock()
	}
	return
}
