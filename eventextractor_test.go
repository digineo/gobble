package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	guesspaths = false
	os.Exit(m.Run())
}

func loadLog(t *testing.T, fname string) []byte {
	t.Helper()

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("cannot read log file %s: %v", fname, err)
	}
	return data
}

// JSON files in testdata/ are kept in a human readable form. This means
// we need to decode and re-encode its contents to generate comparable
// results.
func loadEvent(t *testing.T, fname string) []byte {
	t.Helper()

	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("cannot open event file %s: %v", fname, err)
	}
	defer f.Close()

	var evt sentry.Event
	if err := json.NewDecoder(f).Decode(&evt); err != nil {
		t.Fatalf("cannot parse event file %s: %v", fname, err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&evt); err != nil {
		panic(err) // how?
	}

	return bytes.TrimSpace(buf.Bytes())
}

func Test_extractEvent(t *testing.T) {
	tt := []struct {
		name      string
		logFile   string
		eventJSON string
	}{
		{"simple panic", "simplepanic.raw", "simplepanic.json"},
		{"no panic", "nopanic.raw", ""},
	}
	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			input := loadLog(t, "testdata/"+tc.logFile)

			event, err := extractEvent(bytes.NewBuffer(input))
			require.NoError(t, err)

			if tc.eventJSON == "" {
				assert.Nil(t, event)
				return
			}

			actual, err := json.Marshal(&event)
			require.NoError(t, err)

			expected := loadEvent(t, "testdata/"+tc.eventJSON)
			assert.Equal(t, expected, actual)
		})
	}
}
