package streaming

import (
	"sync"
	"time"
	"waveguide/lib/adn"
)

type Context struct {
	mtx     sync.Mutex
	streams map[string]*StreamInfo
}

func NewContext() *Context {
	return &Context{
		streams: make(map[string]*StreamInfo),
	}
}

func (ctx *Context) Online(limit int) (streams []*StreamInfo) {
	ctx.mtx.Lock()
	for k := range ctx.streams {
		streams = append(streams, ctx.streams[k])
		limit--
		if limit == 0 {
			break
		}
	}
	ctx.mtx.Unlock()
	return
}

func (ctx *Context) Ensure(uid adn.UID, username string, chatid adn.ChanID) (i *StreamInfo) {
	k := uid.String()
	if len(k) > 0 {
		ctx.mtx.Lock()
		var ok bool
		i, ok = ctx.streams[k]
		if !ok {
			i = new(StreamInfo)
			i.ID = uid
			i.Username = username
			i.ChatID = chatid
			i.LastUpdate = time.Now()
			ctx.streams[k] = i
		}
		ctx.mtx.Unlock()
	}
	return
}

func (ctx *Context) Has(k string) (has bool) {
	if len(k) > 0 {
		ctx.mtx.Lock()
		_, has = ctx.streams[k]
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

func (ctx *Context) Remove(k string) {
	if len(k) > 0 {
		ctx.mtx.Lock()
		_, has := ctx.streams[k]
		if has {
			delete(ctx.streams, k)
		}
		ctx.mtx.Unlock()
	}
}
