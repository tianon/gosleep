package progress_test

import (
	"fmt"
	"os"
	"time"

	"go.tianon.xyz/progress"
)

func ExampleBar() {
	bar := progress.NewBar(os.Stdout)
	bar.Max = 100

	bar.Start()
	for bar.Val = bar.Min; bar.Val <= bar.Max; bar.Val += 25 {
		bar.Tick()
		time.Sleep(time.Second)
	}
	bar.Finish()
}

func ExampleBar_tickString() {
	bar := progress.NewBar(nil)
	bar.Min = -100
	bar.Max = 100

	bar.Prefix = func(_ *progress.Bar) string {
		return "["
	}
	bar.Suffix = func(b *progress.Bar) string {
		return fmt.Sprintf("]%3.0f%%", b.Progress()*100)
	}

	bar.Phases = []string{
		" ",
		"-",
		"=",
	}

	for bar.Val = bar.Min; bar.Val <= bar.Max; bar.Val += 25 {
		fmt.Println(bar.TickString(8))
	}

	// Output:
	// [  ]  0%
	// [  ] 12%
	// [- ] 25%
	// [- ] 38%
	// [= ] 50%
	// [= ] 62%
	// [=-] 75%
	// [=-] 88%
	// [==]100%
}

func ExampleBar_tickStringManyPhases() {
	bar := progress.NewBar(nil)
	bar.Min = -100
	bar.Max = 100

	bar.Prefix = func(_ *progress.Bar) string {
		return "["
	}
	bar.Suffix = func(b *progress.Bar) string {
		return fmt.Sprintf("]%3.0f%%", b.Progress()*100)
	}

	bar.Phases = []string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
	}

	for bar.Val = bar.Min; bar.Val <= bar.Max; bar.Val += 5 {
		fmt.Println(bar.TickString(10))
	}

	// Output:
	// [0000]  0%
	// [0000]  2%
	// [1000]  5%
	// [2000]  8%
	// [3000] 10%
	// [4000] 12%
	// [5000] 15%
	// [6000] 18%
	// [7000] 20%
	// [8000] 22%
	// [9000] 25%
	// [9000] 28%
	// [9100] 30%
	// [9200] 32%
	// [9300] 35%
	// [9400] 38%
	// [9500] 40%
	// [9600] 42%
	// [9700] 45%
	// [9800] 48%
	// [9900] 50%
	// [9900] 52%
	// [9910] 55%
	// [9920] 57%
	// [9930] 60%
	// [9940] 62%
	// [9950] 65%
	// [9960] 68%
	// [9970] 70%
	// [9980] 72%
	// [9990] 75%
	// [9990] 78%
	// [9991] 80%
	// [9992] 82%
	// [9993] 85%
	// [9994] 88%
	// [9995] 90%
	// [9996] 92%
	// [9997] 95%
	// [9998] 98%
	// [9999]100%
}
