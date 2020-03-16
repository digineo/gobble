package main

import (
	"bytes"
	"flag"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func main() {
	dsn := flag.String("dsn", "", "sentry dsn")
	serviceName := flag.String("service", "", "service name")
	flag.Parse()

	if *serviceName == "" {
		logrus.Fatal("service name is needed")
	}
	if *dsn == "" {
		logrus.Fatal("dsn is needed")
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: *dsn,
	})
	if err != nil {
		logrus.WithError(err).Fatal("sentry init failed")
	}
	defer sentry.Flush(time.Second * 5)

	var buf bytes.Buffer
	if err := readServiceJournal(*serviceName, &buf); err != nil {
		sentry.CaptureException(err)
		logrus.WithError(err).Fatal("reading journal failed")
	}

	if e := extractEvent(&buf); e != nil {
		sentry.CaptureEvent(e)
	}
}
