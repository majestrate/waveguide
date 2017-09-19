package util

import "bytes"

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) Close() error {
	b.Reset()
	return nil
}
