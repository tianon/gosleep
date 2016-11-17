package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	jinzhuNow "github.com/jinzhu/now"
	"gopkg.in/cheggaaa/pb.v1"
)

type sleepFlags struct {
	For   bool `long:"for"`
	Until bool `long:"until"`
}

func main() {
	start := time.Now()

	opts := sleepFlags{}
	flagParser := flags.NewParser(&opts, flags.Default)
	flagParser.Usage = "--for <duration> OR --until <time>"

	args, err := flagParser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				return
			}
		}
		os.Exit(1)
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "error: argument expected\n\n")
		flagParser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if opts.For == opts.Until {
		fmt.Fprintf(os.Stderr, "error: exactly one of --for or --until expected\n\n")
		flagParser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	var until time.Time
	if opts.For {
		d, err := time.ParseDuration(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse duration %q: %v\n", args[0], err)
			os.Exit(1)
		}
		until = start.Add(d)
	} else { // opts.Until
		until, err = jinzhuNow.New(start).Parse(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse time %q: %v\n", args[0], err)
			os.Exit(1)
		}
	}

	// TODO consider making these optional flags
	increment := 500 * time.Millisecond
	round := time.Second

	start = start.Round(round)
	until = until.Round(round)

	if until.Before(start) {
		fmt.Fprintf(os.Stderr, "error: requested sleep time in the past: %s\n", until.Sub(start))
		os.Exit(1)
	}

	bar := pb.New64(until.Sub(start).Nanoseconds())
	bar.ShowCounters = true
	bar.ShowPercent = true
	bar.ShowTimeLeft = false
	bar.ShowFinalTime = false
	bar.SetUnits(pb.U_DURATION)

	bar.Start()
	for now := time.Now(); now.Before(until); now = time.Now() {
		bar.Set64(now.Round(round).Sub(start).Nanoseconds())
		time.Sleep(increment)
	}
	bar.Finish()
}
