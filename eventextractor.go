package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/maruel/panicparse/stack"
)

var guesspaths = true // false in tests

func extractEvent(r io.Reader) (*sentry.Event, error) {
	var panicBuf PanicFilter
	c, err := stack.ParseDump(r, &panicBuf, guesspaths)
	if err != nil {
		return nil, fmt.Errorf("parsing event dump failed: %w", err)
	}

	if c == nil || len(c.Goroutines) == 0 {
		return nil, nil
	}

	evt := sentry.NewEvent()
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
		evt.Threads = append(evt.Threads, t)
	}
	evt.Message = panicBuf.Value()
	evt.Level = sentry.LevelFatal

	return evt, nil
}
