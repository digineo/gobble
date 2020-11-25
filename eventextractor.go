package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/maruel/panicparse/v2/stack"
)

func extractEvent(r io.Reader) (*sentry.Event, error) {
	var panicBuf PanicFilter

	c, _, err := stack.ScanSnapshot(r, &panicBuf, stack.DefaultOpts())
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("parsing event dump failed: %w", err)
	}

	if c == nil || len(c.Goroutines) == 0 {
		return nil, nil
	}

	evt := sentry.NewEvent()
	evt.Message = panicBuf.Value()
	evt.Level = sentry.LevelFatal

	for _, routine := range c.Goroutines {
		var frames []sentry.Frame
		for _, line := range routine.Stack.Calls {
			frames = append(frames, sentry.Frame{
				Function: line.Func.Name,
				Package:  line.Func.ImportPath,
				Filename: line.SrcName,
				AbsPath:  line.LocalSrcPath,
				Lineno:   line.Line,
				InApp:    line.Location != stack.Stdlib,
			})
		}

		t := sentry.Thread{
			ID:         strconv.Itoa(routine.ID),
			Name:       fmt.Sprintf("goroutine-%d", routine.ID),
			Stacktrace: &sentry.Stacktrace{Frames: frames},
			Crashed:    routine.First,
			Current:    routine.State == "running",
		}
		evt.Threads = append(evt.Threads, t)
	}

	return evt, nil
}
