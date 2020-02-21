// Package terminal provides helpers to manipulate the terminal
// This assumes that the terminal is vt100 compatible.
package terminal

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// Terminal width defaults to 80 columns
var termcols int

// init runs at startup and saves the terminal width
func init() {
	termcols = 80
	updateWinSize()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	go func() {
		for {
			// Wait for signal
			<-ch
			updateWinSize()
		}
	}()
}

func updateWinSize() {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err == nil {
		termcols = int(ws.Col)
	}
}

// Color is the range of colors used by the terminal
type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// csi returns the Command Sequence Indicator (ESC + [) followed by the command
func csi(command string) string {
	return "\033[" + command
}

// Foreground sets the foreground to the given color
func Foreground(c Color) string {
	cmd := fmt.Sprintf("3%dm", c)
	return csi(cmd)
}

// Background sets the background to the given color
func Background(c Color) string {
	cmd := fmt.Sprintf("4%dm", c)
	return csi(cmd)
}

// Reset sends the sequence to reset the terminal text to defaults
func Reset() string {
	return csi("m")
}

// Clear sends the sequence to clear the screen
func Clear() string {
	return csi("2J") + csi("H")
}

// Bold makes the text bold
func Bold() string {
	return csi("1m")
}

// Italic makes the text italic
func Italic() string {
	return csi("2m")
}

// Underline makes the text underlined
func Underline() string {
	return csi("3m")
}

// HorizontalLine prints a horizontal line spanning the width of the terminal
func HorizontalLine() string {
	return strings.Repeat("-", termcols)
}

// Display writes the given string to the terminal, wrapping only on word
// boundaries, and only when an additional word would cause the line length
// to exceed the terminal width.
func Display(_ string) {
}
