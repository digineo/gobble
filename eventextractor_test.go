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

func loadTestCase(t *testing.T, name string) (testCase struct {
	raw, event []byte
	hasEvent   bool
}) {
	t.Helper()

	rawData, err := ioutil.ReadFile("testdata/" + name + ".raw")
	if err != nil {
		t.Fatalf("cannot read raw input data for %s: %v", name, err)
	}
	testCase.raw = rawData

	evtData, err := ioutil.ReadFile("testdata/" + name + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return // "nopanic" case
		}
		t.Fatalf("cannot read event JSON for %s: %v", name, err)
	}
	testCase.hasEvent = true

	// JSON files in testdata/ are kept in a human readable form. This means
	// we need to decode and re-encode its contents to generate comparable
	// results.

	var evt sentry.Event
	if err := json.Unmarshal(evtData, &evt); err != nil {
		t.Fatalf("cannot parse event file %s: %v", name, err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&evt); err != nil {
		panic(err) // how?
	}

	testCase.event = bytes.TrimSpace(buf.Bytes())

	return testCase
}

func Test_extractEvent(t *testing.T) {
	tt := []string{
		"nopanic",
		"simplepanic",
		"recovered",
		"concurrent",
	}

	for i := range tt {
		name := tt[i]

		t.Run(name, func(t *testing.T) {
			tc := loadTestCase(t, name)

			actual, err := extractEvent(bytes.NewBuffer(tc.raw))
			require.NoError(t, err)

			if !tc.hasEvent {
				assert.Nil(t, actual)

				return
			}

			actualJSON, err := json.Marshal(actual)
			require.NoError(t, err)

			assert.EqualValues(t, string(tc.event), string(actualJSON))
		})
	}
}
