package terminal

import (
	"fmt"
	"strings"

	"github.com/buger/goterm"
	"github.com/mattn/go-runewidth"
)

func Print(text string) {
	width := goterm.Width()
	height := goterm.Height()

	lines := strings.Split(text, "\n")
	textHeight := len(lines)

	verticalPadding := (height - textHeight - 15) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	for i := 0; i < verticalPadding; i++ {
		fmt.Println()
	}

	for _, line := range lines {
		textWidth := runewidth.StringWidth(line)
		horizontalPadding := (width - textWidth) / 2
		if horizontalPadding < 0 {
			horizontalPadding = 0
		}
		fmt.Println(strings.Repeat(" ", horizontalPadding) + line)
	}
}
