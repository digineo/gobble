package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

func getServiceProperty(service, name string) string {
	var out bytes.Buffer
	c := exec.Command("systemctl", "show", fmt.Sprintf("--property=%s", name), service)
	c.Stdout = &out
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	s := strings.TrimSpace(out.String())
	if i := strings.IndexByte(s, '='); i > 0 {
		return s[i+1:]
	}

	return ""
}

func addMatches(j *sdjournal.Journal, matches []sdjournal.Match) (err error) {
	for _, m := range matches {
		if err = j.AddMatch(m.String()); err != nil {
			return
		}
	}

	return
}

func readServiceJournal(serviceName string, buf io.StringWriter) error {
	timestamp, err := strconv.Atoi(getServiceProperty(serviceName, "ExecMainStartTimestampMonotonic"))
	if err != nil {
		return fmt.Errorf("failed to get timestamp: %w", err)
	}

	j, err := sdjournal.NewJournal()
	if err != nil {
		return fmt.Errorf("cannot open systemd journal: %w", err)
	}
	defer j.Close()

	err = addMatches(j, []sdjournal.Match{{
		Field: sdjournal.SD_JOURNAL_FIELD_PID,
		Value: getServiceProperty(serviceName, "ExecMainPID"),
	}, {
		Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
		Value: serviceName,
	}})
	if err != nil {
		return err
	}

	if err := j.SeekHead(); err != nil {
		return fmt.Errorf("cannot seek to begin of journal: %w", err)
	}

	for {
		c, err := j.Next()
		if err != nil {
			return fmt.Errorf("fetching journal entry failed: %w", err)
		}
		// Return when on the end of journal
		if c == 0 {
			return nil
		}

		entry, err := j.GetEntry()
		if err != nil {
			return fmt.Errorf("reading journal entry failed: %w", err)
		}

		// Filter manually here since  sd_journal_seek_monotonic_usec is not provided by sdjournal.
		if int(entry.MonotonicTimestamp) < timestamp {
			continue
		}

		buf.WriteString(entry.Fields["MESSAGE"] + "\n")
	}
}
