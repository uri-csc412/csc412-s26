package asciiart

import (
	"fmt"
	"time"

	"inklings/internal/sim"
)

// ANSI text colors (same codes as the original C++ TextColor enum)

type TextColor int

const (
	TextDefault TextColor = 39
	TextBlack   TextColor = 30
	TextRed     TextColor = 31
	TextGreen   TextColor = 32
	TextYellow  TextColor = 33
	TextBlue    TextColor = 34
	TextMagenta TextColor = 35
	TextCyan    TextColor = 36
	TextWhite   TextColor = 37
)

func setTextColor(c TextColor) {
	fmt.Printf("\033[%dm", int(c))
}

func resetTextColor() {
	setTextColor(TextDefault)
}

func clearTerminal() {
	fmt.Print("\033[H\033[J")
}

func DrawGridAndInklings(s *sim.State) {
	// TODO: draw ASCII grid and inklings based on s.Grid and s.Inklings.
	_ = s
}

func DrawState(s *sim.State) {
	// TODO: display ink levels, number of live goroutines, etc.
	_ = s
}

func EventLoop(s *sim.State) {
	// TODO: periodically redraw terminal and listen for input / commands.
	for {
		clearTerminal()
		DrawGridAndInklings(s)
		DrawState(s)
		time.Sleep(time.Second)
	}
}
