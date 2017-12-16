package streaming

type Context struct {
	streams map[int64]*StreamInfo
}

func NewContext() *Context {
	return &Context{
		streams: make(map[int64]*StreamInfo),
	}
}

func (ctx *Context) Ensure(k int64) (i *StreamInfo) {
	var ok bool
	i, ok = ctx.streams[k]
	if !ok {
		i = new(StreamInfo)
		ctx.streams[k] = i
	}
	return
}

func (ctx *Context) Find(k int64) (i *StreamInfo) {
	i, _ = ctx.streams[k]
	return
}
