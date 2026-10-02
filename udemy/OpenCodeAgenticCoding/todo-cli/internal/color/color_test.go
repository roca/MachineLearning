package color

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectRejectsNonFileWriter(t *testing.T) {
	t.Parallel()

	if Detect(&bytes.Buffer{}) {
		t.Error("Detect(buffer) = true, want false")
	}
}

func TestDetectRejectsRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	t.Cleanup(func() { f.Close() })

	if Detect(f) {
		t.Error("Detect(regular file) = true, want false")
	}
}

func TestDetectHonoursNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if Detect(os.Stdout) {
		t.Error("Detect with NO_COLOR set = true, want false")
	}
}

func TestDetectHonoursDumbTerminal(t *testing.T) {
	t.Setenv("TERM", "dumb")
	if Detect(os.Stdout) {
		t.Error("Detect with TERM=dumb = true, want false")
	}
}

func TestModeSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want Mode
	}{
		{in: "", want: Auto},
		{in: "auto", want: Auto},
		{in: "AUTO", want: Auto},
		{in: "always", want: Always},
		{in: "true", want: Always},
		{in: "on", want: Always},
		{in: "1", want: Always},
		{in: "never", want: Never},
		{in: "false", want: Never},
		{in: "off", want: Never},
		{in: "0", want: Never},
	}

	for _, tt := range tests {
		var m Mode
		if err := m.Set(tt.in); err != nil {
			t.Errorf("Set(%q): %v", tt.in, err)
			continue
		}
		if m != tt.want {
			t.Errorf("Set(%q) = %v, want %v", tt.in, m, tt.want)
		}
	}
}

func TestModeSetRejectsGarbage(t *testing.T) {
	t.Parallel()

	var m Mode
	if err := m.Set("chartreuse"); err == nil {
		t.Fatal("Set(chartreuse) succeeded, want error")
	}
}

func TestModeString(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		mode Mode
		want string
	}{{Auto, "auto"}, {Always, "always"}, {Never, "never"}} {
		if got := tt.mode.String(); got != tt.want {
			t.Errorf("Mode(%d).String() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestModeIsBoolFlag(t *testing.T) {
	t.Parallel()

	if !Auto.IsBoolFlag() {
		t.Error("IsBoolFlag() = false, want true so -color works without a value")
	}
}

// Detect is false for a pipe, so only Always overrides it in tests.
func TestModeEnabled(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	for _, tt := range []struct {
		mode Mode
		want bool
	}{{Auto, false}, {Always, true}, {Never, false}} {
		if got := tt.mode.Enabled(&buf); got != tt.want {
			t.Errorf("Mode(%v).Enabled(buffer) = %v, want %v", tt.mode, got, tt.want)
		}
	}
}

func TestNewDisabled(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Never)
	if p.Enabled() {
		t.Error("Enabled() = true, want false")
	}
	if got := p.OK("done"); got != "done" {
		t.Errorf("OK(%q) = %q, want unchanged", "done", got)
	}
}

func TestNewForceEnables(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Always)
	if !p.Enabled() {
		t.Fatal("Enabled() = false, want true")
	}
	if got := p.OK("done"); !strings.Contains(got, "done") || Strip(got) != "done" {
		t.Errorf("OK(%q) = %q, want styled %q", "done", got, "done")
	}
}

func TestStylesWrapAndStrip(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Always)
	tests := []struct {
		name  string
		style func(string) string
		text  string
	}{
		{name: "bold", style: p.Bold, text: "TITLE"},
		{name: "muted", style: p.Muted, text: "2026-10-02"},
		{name: "accent", style: p.Accent, text: "7"},
		{name: "warn", style: p.Warn, text: "pending"},
		{name: "ok", style: p.OK, text: "done"},
		{name: "struck", style: p.Struck, text: "buy milk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.style(tt.text)
			if got == tt.text {
				t.Errorf("style(%q) = %q, want escape sequences", tt.text, got)
			}
			if stripped := Strip(got); stripped != tt.text {
				t.Errorf("Strip(style(%q)) = %q, want %q", tt.text, stripped, tt.text)
			}
		})
	}
}

func TestPlainIsUnstyled(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Always)
	if got := p.Plain("buy milk"); got != "buy milk" {
		t.Errorf("Plain(%q) = %q, want unchanged", "buy milk", got)
	}
}

func TestStylesIgnoreEmptyText(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Always)
	if got := p.Warn(""); got != "" {
		t.Errorf("Warn(%q) = %q, want empty", "", got)
	}
}

func TestStrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "no codes", in: "plain text", want: "plain text"},
		{name: "codes", in: "\x1b[1mID\x1b[0m  7", want: "ID  7"},
		{name: "combined", in: "\x1b[2m\x1b[9mbread\x1b[0m", want: "bread"},
		{name: "truncated", in: "abc\x1b[3", want: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Strip(tt.in); got != tt.want {
				t.Errorf("Strip(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEveryStyleEmitsEscapes(t *testing.T) {
	t.Parallel()

	p := New(&bytes.Buffer{}, Always)
	for name, style := range map[string]func(string) string{
		"bold": p.Bold, "muted": p.Muted, "accent": p.Accent,
		"warn": p.Warn, "ok": p.OK, "struck": p.Struck,
	} {
		if got := style("x"); !strings.HasPrefix(got, "\x1b[") {
			t.Errorf("style %s produced no escape sequence: %q", name, got)
		}
	}
}
