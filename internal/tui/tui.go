package tui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsInteractive returns true if stdin is a terminal.
func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}
