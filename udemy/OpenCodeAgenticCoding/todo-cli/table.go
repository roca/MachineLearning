package main

import (
	"io"
	"strings"
	"unicode/utf8"

	"todo-cli/internal/color"
	"todo-cli/internal/todo"
)

// visibleLen is the width of s in terminal columns, ignoring escape sequences.
// Runes count as one column, which is wrong only for double-width characters.
func visibleLen(s string) int {
	return utf8.RuneCountInString(color.Strip(s))
}

const (
	dateLayout = "2006-01-02"
	colGap     = 2
	barIndent  = "\u2502  "      // │ plus two spaces
	tee        = "\u251c\u2500 " // ├─
	elbow      = "\u2514\u2500 " // └─
)

// branch renders the tree prefix for a node: one indent slot per ancestor,
// then a connector for anything below the top level.
func branch(n todo.Node) string {
	var b strings.Builder
	for _, bar := range n.Trail {
		if bar {
			b.WriteString(barIndent)
		} else {
			b.WriteString("   ")
		}
	}
	if n.Depth > 0 {
		if n.Last {
			b.WriteString(elbow)
		} else {
			b.WriteString(tee)
		}
	}
	return b.String()
}

// cell is one table value plus the style used to render it. text should hold
// the unstyled content: widths are measured from it, with any escape sequences
// already present stripped out.
type cell struct {
	text  string
	style func(string) string
}

// column describes a table column's header and alignment.
type column struct {
	head  string
	right bool
}

// writeTable renders rows as aligned, optionally colored columns. Column
// widths are computed from unstyled text so escape sequences never skew
// alignment.
func writeTable(w io.Writer, p color.Palette, cols []column, rows [][]cell) error {
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = visibleLen(c.head)
	}
	for _, row := range rows {
		for i, c := range row {
			widths[i] = max(widths[i], visibleLen(c.text))
		}
	}

	header := make([]cell, len(cols))
	for i, c := range cols {
		header[i] = cell{text: c.head, style: p.Bold}
	}

	var b strings.Builder
	b.WriteString(renderRow(header, cols, widths))
	for _, row := range rows {
		b.WriteByte('\n')
		b.WriteString(renderRow(row, cols, widths))
	}
	b.WriteByte('\n')

	_, err := io.WriteString(w, b.String())
	return err
}

func renderRow(cells []cell, cols []column, widths []int) string {
	var b strings.Builder
	for i, c := range cells {
		style := c.style
		if style == nil {
			style = func(s string) string { return s }
		}
		pad := strings.Repeat(" ", max(widths[i]-visibleLen(c.text), 0))
		if cols[i].right {
			b.WriteString(pad)
			b.WriteString(style(c.text))
		} else {
			b.WriteString(style(c.text))
			b.WriteString(pad)
		}
		if i < len(cells)-1 {
			b.WriteString(strings.Repeat(" ", colGap))
		}
	}
	return strings.TrimRight(b.String(), " ")
}
