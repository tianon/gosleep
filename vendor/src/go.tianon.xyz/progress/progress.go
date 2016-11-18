package progress // import "go.tianon.xyz/progress"

import (
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

type Bar struct {
	Val int64
	Min int64
	Max int64

	Prefix func(b *Bar) string
	Suffix func(b *Bar) string

	Phases []string

	out *os.File
}

var DefaultPhases = []string{
	" ",
	"▏",
	"▎",
	"▍",
	"▌",
	"▋",
	"▊",
	"▉",
	"█",
}

func NewBar(out *os.File) *Bar {
	return &Bar{
		Prefix: func(_ *Bar) string { return " [" },
		Suffix: func(_ *Bar) string { return "] " },
		Phases: DefaultPhases,

		out: out,
	}
}

func (b *Bar) Start() {
	// TODO if isatty
	//b.out.Write([]byte("\x1b[?25l")) // hide cursor?
	b.Tick()
}

func (b *Bar) Finish() {
	b.Tick()
	b.out.Write([]byte("\n"))
	// TODO if isatty
	//b.out.Write([]byte("\x1b[?25h")) // show cursor?
}

// percentage, normalized to 0-100%
func (b *Bar) Progress() float64 {
	if b.Min >= b.Max {
		// ignore bad values like cowards
		return 1.0
	}
	if b.Val < b.Min {
		return 0.0
	}
	if b.Val > b.Max {
		return 1.0
	}
	return float64(b.Val - b.Min) / float64(b.Max - b.Min)
}

func (b *Bar) TickString() string {
	progress := b.Progress()

	//prefix := fmt.Sprintf("%d / %d %s", b.Value, b.Total, b.Prefix)
	prefix := b.Prefix(b)
	suffix := b.Suffix(b)

	width := 80
	if terminal.IsTerminal(int(b.out.Fd())) {
		w, _, err := terminal.GetSize(int(b.out.Fd()))
		if err == nil {
			width = w
		}
	}
	width -= utf8.RuneCountInString(prefix) + utf8.RuneCountInString(suffix)

	filled := float64(width) * progress
	nFull := int(filled)
	phase := int((filled - float64(nFull)) * float64(len(b.Phases)))
	nEmpty := width - nFull

	full := ""
	if nFull >= 0 {
		full = strings.Repeat(b.Phases[len(b.Phases)-1], nFull)
	}

	current := ""
	if phase > 0 && phase < len(b.Phases) {
		current = b.Phases[phase]
	}

	nEmpty = nEmpty-utf8.RuneCountInString(current)
	empty := ""
	if nEmpty >= 0 {
		empty = strings.Repeat(b.Phases[0], nEmpty)
	}

	return strings.Join([]string{
		prefix,
		full,
		current,
		empty,
		suffix,
	}, "")
}

func (b *Bar) Tick() {
	writeln(b.out, b.TickString())
}

// https://github.com/verigak/progress/blob/c5043685c57294129f654c4b736fe5a119b14ec9/progress/helpers.py#L61-L69

func clearln(out *os.File) {
	if terminal.IsTerminal(int(out.Fd())) {
		out.Write([]byte("\r\x1b[K"))
	}
}

func writeln(out *os.File, line string) {
	clearln(out)
	out.Write([]byte(line + "\r"))
}
