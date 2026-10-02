// Command todo is a small command-line todo list manager.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"todo-cli/internal/color"
	"todo-cli/internal/todo"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(os.Stderr, "todo:", err)
		}
		os.Exit(1)
	}
}

type env struct {
	stdout io.Writer
	stderr io.Writer
	now    func() time.Time
}

type command struct {
	name    string
	aliases []string
	summary string
	run     func(e *env, args []string) error
}

var commands = []command{
	{name: "add", summary: "add a todo", run: runAdd},
	{name: "list", aliases: []string{"ls"}, summary: "list todos", run: runList},
	{name: "done", summary: "mark todos as done", run: runDone},
	{name: "reopen", summary: "mark done todos as pending", run: runReopen},
	{name: "rm", aliases: []string{"remove"}, summary: "delete todos", run: runRemove},
	{name: "clear", summary: "delete completed todos", run: runClear},
}

func run(args []string, stdout, stderr io.Writer) error {
	e := &env{stdout: stdout, stderr: stderr, now: time.Now}

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		usage(stdout)
		if len(args) == 0 {
			return errors.New("no command given")
		}
		return nil
	}

	name, rest := args[0], args[1:]
	if name == "help" || name == "h" {
		usage(stdout)
		return nil
	}
	for _, c := range commands {
		if c.name == name || slices.Contains(c.aliases, name) {
			return c.run(e, rest)
		}
	}
	usage(stderr)
	return fmt.Errorf("unknown command %q", name)
}

func runAdd(e *env, args []string) error {
	fs, file := e.flagSet("add", "usage: todo add [-file path] [-parent id] <title>...")
	parent := fs.Int("parent", 0, "id of the parent todo (0 for top level)")
	fs.IntVar(parent, "p", 0, "shorthand for -parent")
	if err := fs.Parse(args); err != nil {
		return err
	}

	title := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if title == "" {
		return errors.New("add requires a title")
	}

	l, s, err := load(*file)
	if err != nil {
		return err
	}
	t, err := l.Add(title, *parent, e.now())
	if err != nil {
		return err
	}
	if err := s.Save(l); err != nil {
		return err
	}
	fmt.Fprintf(e.stdout, "added #%d %s\n", t.ID, t.Title)
	return nil
}

func runList(e *env, args []string) error {
	fs, file := e.flagSet("list", "usage: todo list [-file path] [-all] [-color=never]")
	all := fs.Bool("all", false, "include completed todos")
	mode := color.Auto
	fs.Var(&mode, "color", "colorize output: auto, always, never")
	if err := fs.Parse(args); err != nil {
		return err
	}

	l, _, err := load(*file)
	if err != nil {
		return err
	}

	items := l.Tree(*all)
	if len(items) == 0 {
		fmt.Fprintln(e.stdout, "no todos")
		return nil
	}

	p := color.New(e.stdout, mode)
	rows := make([][]cell, 0, len(items))
	for _, n := range items {
		status, statusStyle := "pending", p.Warn
		title := p.Plain(n.Title)
		if n.Done {
			status, statusStyle = "done", p.OK
			title = p.Struck(n.Title)
		}
		rows = append(rows, []cell{
			{text: strconv.Itoa(n.ID), style: p.Accent},
			{text: p.Muted(branch(n)) + title, style: p.Plain},
			{text: status, style: statusStyle},
			{text: n.CreatedAt.Format(dateLayout), style: p.Muted},
		})
	}
	return writeTable(e.stdout, p, []column{
		{head: "ID", right: true},
		{head: "TITLE"},
		{head: "STATUS"},
		{head: "CREATED"},
	}, rows)
}

func runDone(e *env, args []string) error {
	return e.setDone(args, "done", "usage: todo done [-file path] <id>...")
}

func runReopen(e *env, args []string) error {
	return e.setDone(args, "reopen", "usage: todo reopen [-file path] <id>...")
}

func (e *env) setDone(args []string, name, usage string) error {
	fs, file := e.flagSet(name, usage)
	if err := fs.Parse(args); err != nil {
		return err
	}

	ids, err := todo.ParseIDs(fs.Args())
	if err != nil {
		return err
	}

	l, s, err := load(*file)
	if err != nil {
		return err
	}
	if name == "done" {
		err = l.Complete(e.now(), ids...)
	} else {
		err = l.Reopen(ids...)
	}
	if err != nil {
		return err
	}
	if err := s.Save(l); err != nil {
		return err
	}
	fmt.Fprintf(e.stdout, "%s %s\n", name, plural(len(ids), "todo", "todos"))
	return nil
}

func runRemove(e *env, args []string) error {
	fs, file := e.flagSet("rm", "usage: todo rm [-file path] <id>...")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ids, err := todo.ParseIDs(fs.Args())
	if err != nil {
		return err
	}

	l, s, err := load(*file)
	if err != nil {
		return err
	}
	n, err := l.Remove(ids...)
	if err != nil {
		return err
	}
	if err := s.Save(l); err != nil {
		return err
	}
	fmt.Fprintf(e.stdout, "removed %s\n", plural(n, "todo", "todos"))
	return nil
}

func runClear(e *env, args []string) error {
	fs, file := e.flagSet("clear", "usage: todo clear [-file path] [-all]")
	all := fs.Bool("all", false, "delete pending todos too")
	if err := fs.Parse(args); err != nil {
		return err
	}

	l, s, err := load(*file)
	if err != nil {
		return err
	}

	var n int
	if *all {
		n = l.Clear()
	} else {
		done := l.Filter(true)
		if len(done) > 0 {
			// Prune, not Remove: unfinished subtasks survive and move up a level.
			n, err = l.Prune(idsOf(done)...)
			if err != nil {
				return err
			}
		}
	}
	if err := s.Save(l); err != nil {
		return err
	}
	fmt.Fprintf(e.stdout, "cleared %s\n", plural(n, "todo", "todos"))
	return nil
}

func (e *env) flagSet(name, usage string) (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(e.stderr)
	fs.Usage = func() { fmt.Fprintln(e.stderr, usage) }
	path, err := todo.DefaultPath()
	if err != nil {
		path = "" // reported by load
	}
	file := fs.String("file", path, "path to the todo data file")
	return fs, file
}

func load(path string) (*todo.List, *todo.Store, error) {
	if path == "" {
		return nil, nil, errors.New("no data file: pass -file or set TODO_FILE")
	}
	s := todo.NewStore(path)
	l, err := s.Load()
	if err != nil {
		return nil, nil, err
	}
	return l, s, nil
}

func idsOf(items []todo.Todo) []int {
	ids := make([]int, 0, len(items))
	for _, t := range items {
		ids = append(ids, t.ID)
	}
	return ids
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

func usage(w io.Writer) {
	fmt.Fprint(w, "todo - a small command-line todo list\n\n")
	fmt.Fprint(w, "usage: todo <command> [flags] [args]\n\n")
	fmt.Fprint(w, "commands:\n")
	for _, c := range commands {
		names := strings.Join(append([]string{c.name}, c.aliases...), ", ")
		fmt.Fprintf(w, "  %-16s %s\n", names, c.summary)
	}
	fmt.Fprintf(w, "  %-16s %s\n", "help, h", "show this help")
	fmt.Fprint(w, "\nflags:\n")
	fmt.Fprint(w, "  -file path      todo data file (default $TODO_FILE or <config dir>/todo-cli/todos.json)\n")
	fmt.Fprint(w, "  -color          list colors: auto (default), always, never; also honours NO_COLOR and TERM=dumb\n")
	fmt.Fprint(w, "  -parent id      add a todo as a subtask of id\n")
	fmt.Fprint(w, "\nexamples:\n")
	fmt.Fprint(w, "  todo add \"ship v1.0\"\n")
	fmt.Fprint(w, "  todo add -parent 1 \"write tests\"    # nested under #1\n")
	fmt.Fprint(w, "  todo list -all                      # full tree, colors on a terminal\n")
}
