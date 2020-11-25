package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// Set via `-ldflags`.
var (
	version = "development"
	commit  = ""
)

// Command line flags.
var (
	dsn   = os.Getenv("SENTRY_DSN")
	svc   = ""
	env   = "production"
	debug = false
)

// Use logf for logging, do not use log.Fatal and friends. Those invoke
// syscall.Exit which terminates the program immediately and skips deferred
// functions.
func logf(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func run() bool {
	err := sentry.Init(sentry.ClientOptions{
		Environment: env,
		Dsn:         dsn,
	})
	if err != nil {
		logf("initializing sentry failed: %v", err)

		return false
	}
	defer sentry.Flush(5 * time.Second)

	var buf bytes.Buffer
	if err := readServiceJournal(svc, &buf); err != nil {
		sentry.CaptureException(err)
		logf("reading journal failed: %v", err)

		return false
	}

	if buf.Len() == 0 {
		logf("no event found, empty buffer")

		return true
	}

	evt, err := extractEvent(&buf)
	if err != nil {
		sentry.CaptureException(err)
		logf("extracting event failed: %v", err)

		return true
	}
	if evt == nil {
		logf("no event found")

		return true
	}

	if debug {
		data, _ := json.Marshal(evt)
		logf("sending event: %s", data)
	}

	if id := sentry.CaptureEvent(evt); id != nil {
		logf("event sent, id: %s", *id)
	} else {
		logf("event was not sent")
	}

	return true
}

func main() {
	log.SetFlags(0) // we're running as systemd hook; the time is already present
	logf("gobble version %s (%s)", version, commit)

	flag.StringVar(&dsn, "dsn", dsn, "sentry `DSN`")
	flag.StringVar(&svc, "service", svc, "service `name`")
	flag.StringVar(&env, "env", env, "environment `name`")
	flag.BoolVar(&debug, "debug", debug, "print sentry events to stdout for debugging")
	flag.Parse()

	if svc == "" {
		logf("missing or empty -service flag")
	}
	if dsn == "" {
		logf("missing or empty -dsn flag or SENTRY_DSN environment variable")
	}

	if debug {
		logf("debug: service=%s env=%s dsn=%s", svc, env, dsn)
	}

	if !run() {
		os.Exit(1)
	}
}
