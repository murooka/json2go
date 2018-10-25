package main

import (
	"bytes"
	"fmt"
)

type ExtBuffer struct {
	*bytes.Buffer
}

func NewExtBuffer() *ExtBuffer {
	return &ExtBuffer{
		&bytes.Buffer{},
	}
}

func (b *ExtBuffer) Print(s string) {
	b.WriteString(s)
}

func (b *ExtBuffer) Println(s string) {
	b.WriteString(s)
	b.WriteByte('\n')
}

func (b *ExtBuffer) Printf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format, args...))
}

func (b *ExtBuffer) Printlnf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format, args...))
	b.WriteByte('\n')
}
