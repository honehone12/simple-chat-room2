package main

type Buffer struct {
	inner  [InputBufferSize]byte
	cursor byte
}

func NewBuffer() *Buffer {
	b := Buffer{
		inner:  [InputBufferSize]byte{},
		cursor: 0,
	}
	for i := 0; i < InputBufferSize; i++ {
		b.inner[i] = Space
	}
	return &b
}

func (b *Buffer) Add(value byte) {
	if b.cursor < InputBufferSize {
		b.inner[b.cursor] = value
		b.cursor++
	}
}

func (b *Buffer) Back() {
	if b.cursor > 0 {
		b.cursor--
		b.inner[b.cursor] = Space
	}
}
