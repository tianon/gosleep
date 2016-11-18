package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	jinzhuNow "github.com/jinzhu/now"
	"go.tianon.xyz/progress"
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

	bar := progress.NewBar(os.Stdout)
	bar.Min = start.Unix()
	bar.Max = until.Unix()

	now := start
	bar.Val = now.Unix()
	bar.Prefix = func(b *progress.Bar) string {
		return fmt.Sprintf(" %s / %s [", now.Sub(start).String(), until.Sub(start).String())
	}
	bar.Suffix = func(b *progress.Bar) string {
		return fmt.Sprintf("] % 3.01f%% ", b.Progress()*100)
	}

	bar.Start()
	for now = start; now.Before(until); now = time.Now().Round(round) {
		bar.Val = now.Unix()
		bar.Tick()
		time.Sleep(increment)
	}
	bar.Val = now.Unix()
	bar.Finish()
}
