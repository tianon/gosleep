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
	// TODO  consider making "start" an optional flag
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
	interval := 100 * time.Millisecond
	round := time.Second

	if until.Before(start) {
		fmt.Fprintf(os.Stderr, "error: requested sleep time in the past: %s\n", until.Sub(start))
		os.Exit(1)
	}

	bar := progress.NewBar(os.Stdout)
	bar.Min = start.UnixNano()
	bar.Max = until.UnixNano()

	now := start
	bar.Val = now.UnixNano()
	bar.Prefix = func(b *progress.Bar) string {
		if b.TickWidth() < 50 {
			// if we're really tight, keep it simple
			return " ["
		}

		rNow := now.Round(round)
		rStart := start.Round(round)
		rUntil := until.Round(round)

		return fmt.Sprintf(" %s / %s [", rNow.Sub(rStart).String(), rUntil.Sub(rStart).String())
	}
	bar.Suffix = func(b *progress.Bar) string {
		rNow := now.Round(round)
		//rStart := start.Round(round)
		rUntil := until.Round(round)

		str := fmt.Sprintf("] %5.01f%% ", b.Progress()*100)
		if b.TickWidth() > 100 {
			// if we're extra wide, let's add a little extra detail
			str += fmt.Sprintf("(%s rem) ", rUntil.Sub(rNow).String())
		}
		return str
	}

	bar.Start()
	for now = start; now.Before(until); now = time.Now() {
		bar.Val = now.UnixNano()
		bar.Tick()
		time.Sleep(interval)
	}
	bar.Val = now.UnixNano()
	bar.Finish()
}
