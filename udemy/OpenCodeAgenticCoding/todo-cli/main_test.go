package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// args builds a command invocation with an explicit data file.
func args(cmd, path string, rest ...string) []string {
	return append([]string{cmd, "-file", path}, rest...)
}

func exec(t *testing.T, argv []string) string {
	t.Helper()

	var out, errOut bytes.Buffer
	if err := run(argv, &out, &errOut); err != nil {
		t.Fatalf("run(%q): %v (stderr: %s)", argv, err, errOut.String())
	}
	return out.String()
}

func execErr(t *testing.T, argv []string) error {
	t.Helper()

	var out, errOut bytes.Buffer
	err := run(argv, &out, &errOut)
	if err == nil {
		t.Fatalf("run(%q) succeeded, want error (stdout: %s)", argv, out.String())
	}
	return err
}

func newFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "todos.json")
}

func wantContains(t *testing.T, got string, wants ...string) {
	t.Helper()

	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func wantNotContains(t *testing.T, got string, unwanted ...string) {
	t.Helper()

	for _, bad := range unwanted {
		if strings.Contains(got, bad) {
			t.Errorf("output %q unexpectedly contains %q", got, bad)
		}
	}
}

func TestLifecycle(t *testing.T) {
	t.Parallel()

	path := newFile(t)

	exec(t, args("add", path, "buy milk"))
	exec(t, args("add", path, "walk dog"))
	wantContains(t, exec(t, args("add", path, "write tests")), "added #3 write tests")

	wantContains(t, exec(t, args("done", path, "1")), "done 1 todo")

	out := exec(t, args("list", path))
	wantContains(t, out, "walk dog", "write tests", "pending")
	wantNotContains(t, out, "buy milk")

	out = exec(t, args("list", path, "-all"))
	wantContains(t, out, "buy milk", "walk dog", "write tests", "done", "pending")

	exec(t, args("reopen", path, "1"))
	wantContains(t, exec(t, args("list", path)), "buy milk")

	wantContains(t, exec(t, args("done", path, "1,2")), "done 2 todos")
	wantContains(t, exec(t, args("clear", path)), "cleared 2 todos")

	out = exec(t, args("list", path, "-all"))
	wantContains(t, out, "write tests")
	wantNotContains(t, out, "buy milk")

	wantContains(t, exec(t, args("rm", path, "3")), "removed 1 todo")
	wantContains(t, exec(t, args("list", path, "-all")), "no todos")
}

func TestIDsSurviveAcrossInvocations(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "a"))
	exec(t, args("add", path, "b"))
	exec(t, args("rm", path, "1"))

	wantContains(t, exec(t, args("add", path, "c")), "added #3 c")
}

func TestClearAll(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "a"))
	exec(t, args("add", path, "b"))

	wantContains(t, exec(t, args("clear", path, "-all")), "cleared 2 todos")
	wantContains(t, exec(t, args("list", path, "-all")), "no todos")
}

func TestErrors(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "a"))

	t.Run("no command", func(t *testing.T) {
		t.Parallel()
		execErr(t, nil)
	})

	t.Run("unknown command", func(t *testing.T) {
		t.Parallel()
		execErr(t, []string{"frobnicate"})
	})

	t.Run("add without title", func(t *testing.T) {
		t.Parallel()
		execErr(t, args("add", path))
	})

	t.Run("done without id", func(t *testing.T) {
		t.Parallel()
		execErr(t, args("done", path))
	})

	t.Run("done unknown id", func(t *testing.T) {
		t.Parallel()
		execErr(t, args("done", path, "42"))
	})

	t.Run("rm unknown id", func(t *testing.T) {
		t.Parallel()
		execErr(t, args("rm", path, "42"))
	})

	t.Run("bad id", func(t *testing.T) {
		t.Parallel()
		execErr(t, args("done", path, "abc"))
	})

	t.Run("corrupt data file", func(t *testing.T) {
		t.Parallel()

		bad := newFile(t)
		if err := os.WriteFile(bad, []byte("{oops"), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		execErr(t, args("list", bad))
	})
}

func TestSubtasksRenderAsTree(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("add", path, "-parent", "2", "delta"))
	exec(t, args("add", path, "-parent", "1", "charlie"))
	exec(t, args("add", path, "foxtrot"))

	want := "ID  TITLE           STATUS   CREATED\n" +
		" 1  alpha           pending  2026-10-02\n" +
		" 2  │  ├─ bravo     pending  2026-10-02\n" +
		" 3  │  │  └─ delta  pending  2026-10-02\n" +
		" 4  │  └─ charlie   pending  2026-10-02\n" +
		" 5  foxtrot         pending  2026-10-02\n"
	if got := exec(t, args("list", path, "-all")); got != want {
		t.Errorf("tree:\n%s\nwant:\n%s", got, want)
	}
}

func TestNestedSubtasksRenderDeeper(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("add", path, "-parent", "2", "delta"))
	exec(t, args("add", path, "-parent", "3", "echo"))
	exec(t, args("add", path, "-parent", "3", "foxtrot"))

	// One root only, so no vertical bars: each level just indents.
	want := "ID  TITLE                STATUS   CREATED\n" +
		" 1  alpha                pending  2026-10-02\n" +
		" 2     └─ bravo          pending  2026-10-02\n" +
		" 3        └─ delta       pending  2026-10-02\n" +
		" 4           ├─ echo     pending  2026-10-02\n" +
		" 5           └─ foxtrot  pending  2026-10-02\n"
	if got := exec(t, args("list", path, "-all")); got != want {
		t.Errorf("tree:\n%s\nwant:\n%s", got, want)
	}
}

func TestAddSubtaskShorthandParent(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-p", "1", "bravo"))

	out := exec(t, args("list", path, "-all"))
	wantContains(t, out, "└─ bravo")
}

func TestAddUnderUnknownParent(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	execErr(t, args("add", path, "-parent", "42", "orphan"))
}

func TestListKeepsDoneParentForPendingChild(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("done", path, "1"))

	out := exec(t, args("list", path)) // pending only
	wantContains(t, out, "alpha", "└─ bravo")
	wantNotContains(t, out, "no todos")
}

func TestListHidesFullyDoneSubtree(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("done", path, "1,2"))

	wantContains(t, exec(t, args("list", path)), "no todos")
}

func TestRemoveParentRemovesSubtree(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("add", path, "-parent", "2", "delta"))
	exec(t, args("add", path, "-parent", "1", "charlie"))

	wantContains(t, exec(t, args("rm", path, "1")), "removed 4 todos")
	wantContains(t, exec(t, args("list", path, "-all")), "no todos")
}

func TestClearKeepsPendingSubtask(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("done", path, "1"))

	wantContains(t, exec(t, args("clear", path)), "cleared 1 todo")

	// bravo survives, promoted to the top level.
	out := exec(t, args("list", path, "-all"))
	wantContains(t, out, "bravo")
	wantNotContains(t, out, "alpha", "└─", "├─", "│")
}

func TestClearKeepsDeepSubtask(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "alpha"))
	exec(t, args("add", path, "-parent", "1", "bravo"))
	exec(t, args("add", path, "-parent", "2", "delta"))
	exec(t, args("done", path, "1,2"))

	wantContains(t, exec(t, args("clear", path)), "cleared 2 todos")
	wantContains(t, exec(t, args("list", path, "-all")), "delta")
}

func TestListColorFlag(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "buy milk"))
	exec(t, args("done", path, "1"))

	var out, errOut bytes.Buffer
	if err := run(args("list", path, "-all"), &out, &errOut); err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.Contains(out.String(), "\x1b[") {
		t.Errorf("piped list output contains escape sequences:\n%q", out.String())
	}

	out.Reset()
	if err := run(args("list", path, "-all", "-color"), &out, &errOut); err != nil {
		t.Fatalf("list -color: %v", err)
	}
	if !strings.Contains(out.String(), "\x1b[") {
		t.Errorf("list -color output has no escape sequences:\n%q", out.String())
	}

	out.Reset()
	if err := run(args("list", path, "-all", "-color=always"), &out, &errOut); err != nil {
		t.Fatalf("list -color=always: %v", err)
	}
	if !strings.Contains(out.String(), "\x1b[") {
		t.Errorf("list -color=always output has no escape sequences:\n%q", out.String())
	}

	for _, mode := range []string{"never", "false", "off", "0"} {
		out.Reset()
		if err := run(args("list", path, "-all", "-color="+mode), &out, &errOut); err != nil {
			t.Fatalf("list -color=%s: %v", mode, err)
		}
		if strings.Contains(out.String(), "\x1b[") {
			t.Errorf("list -color=%s contains escape sequences:\n%q", mode, out.String())
		}
	}

	out.Reset()
	if err := run(args("list", path, "-all", "-color=chartreuse"), &out, &errOut); err == nil {
		t.Error("list -color=chartreuse succeeded, want error")
	}
}

func TestHelp(t *testing.T) {
	t.Parallel()

	wantContains(t, exec(t, []string{"help"}), "usage: todo", "add", "list", "done", "rm", "clear")
	wantContains(t, exec(t, []string{"--help"}), "usage: todo")
}

func TestAliases(t *testing.T) {
	t.Parallel()

	path := newFile(t)
	exec(t, args("add", path, "a"))
	exec(t, args("done", path, "1"))

	wantContains(t, exec(t, args("remove", path, "1")), "removed 1 todo")
	wantContains(t, exec(t, args("ls", path)), "no todos")
}
