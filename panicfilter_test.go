package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicFilter(t *testing.T) {
	tt := []struct {
		name     string
		inp      string
		expected string
	}{
		{"no panic", "This is nonsense\nthis is more nonsense", ""},
		{"single line panic", "This is nonsense\npanic: I died!", "I died!"},
		{"multi line panic", "This is nonsense\npanic: I died!\navenge me!", "I died!\navenge me!"},
	}
	for _, tt := range tt {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			p := &PanicFilter{}
			p.Write([]byte(tt.inp))

			assert.Equal(tt.expected, p.Value())
		})
	}
}
