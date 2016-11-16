package main

import (
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	start := time.Now()

	// TODO accept either duration or absolute time for sleep target
	d, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic(err)
	}

	until := start.Add(d)

	bar := pb.New64(until.Sub(start).Nanoseconds())
	bar.ShowCounters = true
	bar.ShowPercent = true
	bar.ShowTimeLeft = false
	bar.ShowFinalTime = false
	bar.SetUnits(pb.U_DURATION)

	bar.Start()
	for now := time.Now(); now.Before(until); now = time.Now() {
		bar.Set64(now.Round(time.Second).Sub(start.Round(time.Second)).Nanoseconds())
		time.Sleep(500 * time.Millisecond)
	}
	bar.Finish()
}
