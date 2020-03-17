package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/maruel/panicparse/stack"
	"github.com/sirupsen/logrus"
)

func extractEvent(r io.Reader) *sentry.Event {
	var panicBuf PanicFilter
	c, err := stack.ParseDump(r, &panicBuf, true)
	if err != nil {
		sentry.CaptureException(err)
		logrus.WithError(err).Fatal("parsing panic message failed")
	}

	if c == nil || len(c.Goroutines) == 0 {
		return nil
	}

	e := sentry.NewEvent()
	for _, routine := range c.Goroutines {
		var frames []sentry.Frame
		for _, line := range routine.Stack.Calls {
			frames = append(frames, sentry.Frame{
				Function: line.Func.Name(),
				Package:  line.Func.PkgName(),
				Filename: filepath.Base(line.SrcPath),
				AbsPath:  line.LocalSrcPath,
				Lineno:   line.Line,
				InApp:    !line.IsStdlib,
			})
		}

		stacktrace := &sentry.Stacktrace{
			Frames: frames,
		}

		t := sentry.Thread{
			ID:         strconv.Itoa(routine.ID),
			Name:       fmt.Sprintf("goroutine-%d", routine.ID),
			Stacktrace: stacktrace,
			Crashed:    routine.First,
			Current:    routine.State == "running",
		}
		e.Threads = append(e.Threads, t)
	}
	e.Message = panicBuf.Value()
	e.Level = sentry.LevelFatal

	return e
}
