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

type JournalReader struct {
	serviceName string
	isUserService bool
}

func (jr *JournalReader) property(name string) string {
	args := []string{"show"}
	if jr.isUserService {
		args = append(args, "--user")
	}
	args = append(args, fmt.Sprintf("--property=%s", name), jr.serviceName)
	cmd := exec.Command("systemctl", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	s := strings.TrimSpace(out.String())
	if i := strings.IndexByte(s, '='); i > 0 {
		return s[i+1:]
	}

	return ""
}

func (jr *JournalReader) pidMatch() string {
	m := sdjournal.Match{
		Field: sdjournal.SD_JOURNAL_FIELD_PID,
		Value: jr.property("ExecMainPID"),
	}

	return m.String()
}

func (jr *JournalReader) unitMatch() string {
	field := sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT
	if jr.isUserService {
		field = sdjournal.SD_JOURNAL_FIELD_SYSTEMD_USER_UNIT
	}

	m := sdjournal.Match{
		Field: field,
		Value: jr.serviceName,
	}

	return m.String()
}

func (jr *JournalReader) readInto(w io.StringWriter) error {
	timestamp, err := strconv.Atoi(jr.property("ExecMainStartTimestampMonotonic"))
	if err != nil {
		return fmt.Errorf("failed to get timestamp: %w", err)
	}

	j, err := sdjournal.NewJournal()
	if err != nil {
		return fmt.Errorf("cannot open systemd journal: %w", err)
	}
	defer j.Close()

	if err = j.AddMatch(jr.pidMatch()); err != nil {
		return fmt.Errorf("failed to add PID match: %w", err)
	}
	if err = j.AddMatch(jr.unitMatch()); err != nil {
		return fmt.Errorf("failed to add unit match: %w", err)
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

		w.WriteString(entry.Fields["MESSAGE"] + "\n")
	}
}
