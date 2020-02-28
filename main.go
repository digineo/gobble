package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/maruel/panicparse/stack"
)

func main() {
	serviceName := "acrashy.service"

	var buf bytes.Buffer

	readServiceJournal(serviceName, &buf)

	c, err := stack.ParseDump(&buf, os.Stdout, true)
	if err != nil {
		return
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
			if l := len(line.Func.PkgName()); l > pkgLen {
				pkgLen = l
			}
		}
	}

	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}
		if c := bucket.CreatedByString(false); c != "" {
			extra += " [Created by " + c + "]"
		}
		fmt.Printf("%d: %s%s\n", len(bucket.IDs), bucket.State, extra)

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			fmt.Printf(
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.SrcLine(),
				line.Func.Name(), &line.Args)
		}
		if bucket.Stack.Elided {
			io.WriteString(os.Stdout, "    (...)\n")
		}
	}
}
