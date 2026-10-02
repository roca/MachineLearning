package main

import (
	"bytes"
	"strings"
	"testing"

	"todo-cli/internal/color"
)

func render(t *testing.T, mode color.Mode, cols []column, rows [][]cell) string {
	t.Helper()

	var out bytes.Buffer
	p := color.New(&out, mode)
	if err := writeTable(&out, p, cols, rows); err != nil {
		t.Fatalf("writeTable: %v", err)
	}
	return out.String()
}

func TestWriteTableRightAlignsFirstColumn(t *testing.T) {
	t.Parallel()

	p := color.New(&bytes.Buffer{}, color.Never)
	cols := []column{{head: "ID", right: true}, {head: "TITLE"}, {head: "STATUS"}}
	rows := [][]cell{
		{{text: "1", style: p.Accent}, {text: "aa", style: p.Plain}, {text: "pending", style: p.Warn}},
		{{text: "200", style: p.Accent}, {text: "bbb", style: p.Plain}, {text: "done", style: p.OK}},
	}

	want := " ID  TITLE  STATUS\n" +
		"  1  aa     pending\n" +
		"200  bbb    done\n"
	if got := render(t, color.Never, cols, rows); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestWriteTableAlignsMultiWordTitles(t *testing.T) {
	t.Parallel()

	cols := []column{{head: "ID", right: true}, {head: "TITLE"}, {head: "STATUS"}, {head: "CREATED"}}
	rows := [][]cell{
		{{text: "7"}, {text: "buy milk"}, {text: "pending"}, {text: "2026-10-02"}},
		{{text: "12"}, {text: "walk the dog"}, {text: "done"}, {text: "2026-10-02"}},
	}

	want := "ID  TITLE         STATUS   CREATED\n" +
		" 7  buy milk      pending  2026-10-02\n" +
		"12  walk the dog  done     2026-10-02\n"
	if got := render(t, color.Never, cols, rows); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestWriteTableMeasuresWidthInRunes(t *testing.T) {
	t.Parallel()

	cols := []column{{head: "TITLE"}, {head: "STATUS"}}
	rows := [][]cell{
		{{text: "café"}, {text: "pending"}},
		{{text: "naïve longer"}, {text: "done"}},
	}

	want := "TITLE         STATUS\n" +
		"café          pending\n" +
		"naïve longer  done\n"
	if got := render(t, color.Never, cols, rows); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestWriteTableColorsWithoutBreakingLayout(t *testing.T) {
	t.Parallel()

	build := func(p color.Palette) ([]column, [][]cell) {
		return []column{{head: "ID", right: true}, {head: "TITLE"}, {head: "STATUS"}, {head: "CREATED"}},
			[][]cell{
				{{text: "7", style: p.Accent}, {text: "buy milk", style: p.Plain}, {text: "pending", style: p.Warn}, {text: "2026-10-02", style: p.Muted}},
				{{text: "8", style: p.Accent}, {text: "walk the dog", style: p.Struck}, {text: "done", style: p.OK}, {text: "2026-10-02", style: p.Muted}},
			}
	}

	p := color.New(&bytes.Buffer{}, color.Always)
	cols, rows := build(p)
	got := render(t, color.Always, cols, rows)

	plainCols, plainRows := build(color.New(&bytes.Buffer{}, color.Never))
	want := render(t, color.Never, plainCols, plainRows)
	if stripped := color.Strip(got); stripped != want {
		t.Errorf("stripped output:\n%q\nwant:\n%q", stripped, want)
	}

	for _, want := range []string{
		"\x1b[36m7\x1b[0m",                  // id
		"\x1b[33mpending\x1b[0m",            // pending status
		"\x1b[32mdone\x1b[0m",               // done status
		"\x1b[2m\x1b[9mwalk the dog\x1b[0m", // completed title
		"\x1b[2m2026-10-02\x1b[0m",          // date
		"\x1b[1mTITLE\x1b[0m",               // header
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%q", want, got)
		}
	}
}

func TestWriteTablePlainStyleIsUnstyled(t *testing.T) {
	t.Parallel()

	p := color.New(&bytes.Buffer{}, color.Always)
	cols := []column{{head: "TITLE"}}
	rows := [][]cell{{{text: "buy milk", style: p.Plain}}}

	got := render(t, color.Always, cols, rows)
	if want := "\x1b[1mTITLE\x1b[0m\nbuy milk\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteTableTrimsTrailingSpace(t *testing.T) {
	t.Parallel()

	cols := []column{{head: "TITLE"}, {head: "STATUS"}}
	rows := [][]cell{{{text: "a"}, {text: "done"}}}

	for line := range strings.SplitSeq(strings.TrimRight(render(t, color.Never, cols, rows), "\n"), "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line %q ends with a space", line)
		}
	}
}

func TestWriteTableHeaderOnly(t *testing.T) {
	t.Parallel()

	cols := []column{{head: "ID"}, {head: "TITLE"}}
	if got, want := render(t, color.Never, cols, nil), "ID  TITLE\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteTableUnicodeTitles(t *testing.T) {
	t.Parallel()

	cols := []column{{head: "TITLE"}, {head: "STATUS"}}
	rows := [][]cell{{{text: "日本語のタスク"}, {text: "pending"}}}

	if got, want := render(t, color.Never, cols, rows), "TITLE    STATUS\n日本語のタスク  pending\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
