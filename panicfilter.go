package main

import (
	"bytes"
	"strings"
)

const (
	panicHeader = "panic: "
)

type PanicFilter struct {
	found bool
	buf   bytes.Buffer
}

func (p *PanicFilter) Write(d []byte) (n int, err error) {
	return p.WriteString(string(d))
}

func (p *PanicFilter) WriteString(s string) (n int, err error) {
	if !p.found {
		off := strings.Index(s, panicHeader)
		if off < 0 {
			return len(s), nil
		}

		p.found = true
		s = s[off+len(panicHeader):]
	}
	return p.buf.WriteString(s)
}

func (p *PanicFilter) Value() string {
	return p.buf.String()
}
