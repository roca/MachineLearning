// Package color provides minimal ANSI styling for terminal output.
package color

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	reset      = "\x1b[0m"
	bold       = "\x1b[1m"
	dim        = "\x1b[2m"
	yellow     = "\x1b[33m"
	green      = "\x1b[32m"
	cyan       = "\x1b[36m"
	strikethru = "\x1b[9m"
	csi        = "\x1b["
)

// Mode selects when styling is applied.
type Mode int

const (
	// Auto styles output only when the destination is a terminal.
	Auto Mode = iota
	// Always styles output, even when piped.
	Always
	// Never styles output.
	Never
)

// Set implements flag.Value. It accepts the spellings of a boolean flag
// (-color, -color=false) plus the explicit modes "auto", "always", "never".
func (m *Mode) Set(s string) error {
	switch strings.ToLower(s) {
	case "1", "true", "t", "on", "yes", "always":
		*m = Always
	case "0", "false", "f", "off", "no", "never":
		*m = Never
	case "", "auto":
		*m = Auto
	default:
		return fmt.Errorf("invalid color mode %q (want auto, always, or never)", s)
	}
	return nil
}

// String implements flag.Value.
func (m Mode) String() string {
	switch m {
	case Always:
		return "always"
	case Never:
		return "never"
	default:
		return "auto"
	}
}

// IsBoolFlag reports that the flag may be used without a value, so both
// -color and -color=never work.
func (m Mode) IsBoolFlag() bool { return true }

// Enabled reports whether mode styles output written to w.
func (m Mode) Enabled(w io.Writer) bool {
	switch m {
	case Always:
		return true
	case Never:
		return false
	default:
		return Detect(w)
	}
}

// Detect reports whether w is an interactive terminal worth styling. It
// honours the NO_COLOR and TERM=dumb conventions.
func Detect(w io.Writer) bool {
	if _, no := os.LookupEnv("NO_COLOR"); no {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// Palette styles text for one writer. A disabled Palette is a no-op, so
// callers can build output unconditionally.
type Palette struct {
	enabled bool
}

// New returns a Palette for w that follows mode. Use Always to keep styling
// through a pipe, for example into "less -R".
func New(w io.Writer, mode Mode) Palette {
	return Palette{enabled: mode.Enabled(w)}
}

// Enabled reports whether the Palette emits escape sequences.
func (p Palette) Enabled() bool { return p.enabled }

func (p Palette) wrap(codes, s string) string {
	if !p.enabled || s == "" {
		return s
	}
	return codes + s + reset
}

// Bold renders emphasised text.
func (p Palette) Bold(s string) string { return p.wrap(bold, s) }

// Muted renders de-emphasised text.
func (p Palette) Muted(s string) string { return p.wrap(dim, s) }

// Accent renders identifiers.
func (p Palette) Accent(s string) string { return p.wrap(cyan, s) }

// Warn renders pending state.
func (p Palette) Warn(s string) string { return p.wrap(yellow, s) }

// OK renders completed state.
func (p Palette) OK(s string) string { return p.wrap(green, s) }

// Struck renders finished text.
func (p Palette) Struck(s string) string { return p.wrap(dim+strikethru, s) }

// Plain returns s unchanged; it keeps table cells readable without a style.
func (p Palette) Plain(s string) string { return s }

// Strip removes ANSI escape sequences from s.
func Strip(s string) string {
	var b strings.Builder
	for {
		i := strings.Index(s, csi)
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		b.WriteString(s[:i])
		s = s[i+len(csi):]
		j := strings.IndexByte(s, 'm')
		if j < 0 {
			return b.String() // truncated sequence
		}
		s = s[j+1:]
	}
}
