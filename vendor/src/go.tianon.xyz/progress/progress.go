package progress // import "go.tianon.xyz/progress"

import (
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

// Bar is the type representing a standard progress bar.
// It can be used either via TickString (for explicit string generation) or Start/Tick/Finish for direct terminal/os.File output.
type Bar struct {
	Val int64
	Min int64
	Max int64

	// Functions invoked by TickString for generating the text before/after a given progress string.
	Prefix func(b *Bar) string
	Suffix func(b *Bar) string

	// The set of characters used for the progress bar display.
	// The first phase is the implied "empty" and the final phase is the implied "full".
	// See DefaultPhases for the default value provided by NewBar.
	// Note: TickString assumes each of these represents a single output character.
	Phases []string

	Out *os.File
}

// Default value for "Phases" provided by NewBar.
// This default allows for sub-character progress precision.
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

// NewBar creates a new progress bar targeted at "out" (can be "nil" if "TickString" use is intended instead of Start/Tick/Finish).
func NewBar(out *os.File) *Bar {
	return &Bar{
		Prefix: func(_ *Bar) string { return " [" },
		Suffix: func(_ *Bar) string { return "] " },
		Phases: DefaultPhases,

		Out: out,
	}
}

// Start begins progress bar output (invokes Tick).
func (b *Bar) Start() {
	// TODO if isatty
	//b.Out.Write([]byte("\x1b[?25l")) // hide cursor?
	b.Tick()
}

// Finish completes progress bar output (invokes Tick and writes "\n").
func (b *Bar) Finish() {
	b.Tick()
	b.Out.Write([]byte("\n"))
	// TODO if isatty
	//b.Out.Write([]byte("\x1b[?25h")) // show cursor?
}

// Progress returns the current percentage (Val along the line Min <-> Max), normalized to 0-100% as a 0.0-1.0 float64.
//
// The following special cases apply:
//
//   Min >= Max: 100%
//   Val <  Min:   0%
//   Val >  Max: 100%
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
	return float64(b.Val-b.Min) / float64(b.Max-b.Min)
}

// TermWidth returns the width of terminal "out" or -1 if it is not a terminal or if the dimensions of it cannot be determined.
func TermWidth(out *os.File) int {
	if out != nil && terminal.IsTerminal(int(out.Fd())) {
		w, _, err := terminal.GetSize(int(out.Fd()))
		if err == nil {
			return w
		}
	}
	return -1
}

// TickWidth returns the available width for the entire progress bar (or 80 if the available width cannot be determined).
func (b *Bar) TickWidth() int {
	width := TermWidth(b.Out)
	if width < 0 {
		width = 80
	}
	return width
}

// Tick updates progress bar output.
func (b *Bar) Tick() {
	writeln(b.Out, b.TickString(b.TickWidth()))
}

// TickString returns a current progress bar string of "width" (possibly more depending on whether "Prefix" and "Suffix" take all available space).
func (b *Bar) TickString(width int) string {
	prefix := b.Prefix(b)
	suffix := b.Suffix(b)

	width -= utf8.RuneCountInString(prefix) + utf8.RuneCountInString(suffix)

	if width <= 0 {
		// if we already don't have enough space for a progress bar after subtracting Prefix and Suffix, let's force ourselves to get at least a single character
		width = 1
	}

	// https://github.com/verigak/progress/blob/c5043685c57294129f654c4b736fe5a119b14ec9/progress/bar.py#L67-L79

	progress := b.Progress()
	filled := float64(width) * progress
	nFull := int(filled)
	phase := int((filled - float64(nFull)) * float64(len(b.Phases)-1))
	nEmpty := width - nFull

	full := ""
	if nFull >= 0 {
		full = strings.Repeat(b.Phases[len(b.Phases)-1], nFull)
	}

	current := ""
	if phase > 0 && phase < len(b.Phases) {
		current = b.Phases[phase]
	}

	nEmpty = nEmpty - utf8.RuneCountInString(current)
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
