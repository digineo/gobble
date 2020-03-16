package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	simplePanicRaw = `panic: i died, who guessed.
goroutine 1 [running]:
main.iLoveCrashing(...)
	/home/ags/git/digineo.de/acrashyapp/main.go:9
main.innocentOperation(...)
	/home/ags/git/digineo.de/acrashyapp/main.go:13
main.main()
	/home/ags/git/digineo.de/acrashyapp/main.go:25 +0x77`

	simplePanicJSON = `{"level":"fatal","message":"i died, who guessed.\n","sdk":{},"threads":[{"id":"1","name":"routine-1","stacktrace":{"frames":[{"function":"iLoveCrashing","package":"main","filename":"main.go","abs_path":"/usr/lib/go/home/ags/git/digineo.de/acrashyapp/main.go","lineno":9,"in_app":true},{"function":"innocentOperation","package":"main","filename":"main.go","abs_path":"/usr/lib/go/home/ags/git/digineo.de/acrashyapp/main.go","lineno":13,"in_app":true},{"function":"main","package":"main","filename":"main.go","abs_path":"/usr/lib/go/home/ags/git/digineo.de/acrashyapp/main.go","lineno":25,"in_app":true}]},"crashed":true,"current":true}],"user":{},"request":{}}`

	noPanicRaw = `There is no panic here. Shoo shoo.`
)

func Test_extractEvent(t *testing.T) {
	tt := []struct {
		name string
		inp  string
		want string
	}{
		{"simple panic", simplePanicRaw, simplePanicJSON},
		{"no panic", noPanicRaw, ""},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)

			e := extractEvent(bytes.NewBufferString(tc.inp))
			if tc.want == "" {
				assert.Nil(e)
			} else {
				actual, err := json.Marshal(e)
				require.NoError(err)

				assert.Equal(tc.want, string(actual))
			}
		})
	}
}
