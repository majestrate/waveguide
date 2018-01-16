package streaming

type Context struct {
	streams map[string]*StreamInfo
}

func NewContext() *Context {
	return &Context{
		streams: make(map[string]*StreamInfo),
	}
}

func (ctx *Context) Ensure(k string) (i *StreamInfo) {
	if len(k) > 0 {
		var ok bool
		i, ok = ctx.streams[k]
		if !ok {
			i = new(StreamInfo)
			ctx.streams[k] = i
		}
	}
	return
}

func (ctx *Context) Find(k string) (i *StreamInfo) {
	if len(k) > 0 {
		i, _ = ctx.streams[k]
	}
	return
}
