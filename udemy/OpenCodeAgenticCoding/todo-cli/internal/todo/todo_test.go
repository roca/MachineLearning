package todo

import (
	"errors"
	"slices"
	"testing"
	"time"
)

var created = time.Date(2026, time.October, 2, 12, 0, 0, 0, time.UTC)

func newList(titles ...string) *List {
	l := &List{NextID: 1}
	for _, title := range titles {
		if _, err := l.Add(title, 0, created); err != nil {
			panic(err) // parent 0 is always valid
		}
	}
	return l
}

// add adds a todo, failing the test if the list rejects it.
func add(t *testing.T, l *List, title string, parentID int) Todo {
	t.Helper()

	got, err := l.Add(title, parentID, created)
	if err != nil {
		t.Fatalf("Add(%q, parent %d): %v", title, parentID, err)
	}
	return got
}

func ids(items []Todo) []int {
	out := make([]int, 0, len(items))
	for _, t := range items {
		out = append(out, t.ID)
	}
	return out
}

func TestAddAssignsMonotonicIDs(t *testing.T) {
	t.Parallel()

	l := newList("write tests", "ship it")
	if got, want := ids(l.Todos), []int{1, 2}; !slices.Equal(got, want) {
		t.Fatalf("ids = %v, want %v", got, want)
	}
	if l.NextID != 3 {
		t.Errorf("NextID = %d, want 3", l.NextID)
	}
}

func TestAddAfterRemoveKeepsIDsUnique(t *testing.T) {
	t.Parallel()

	l := newList("a", "b", "c")
	if _, err := l.Remove(2); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if got := add(t, l, "d", 0).ID; got != 4 {
		t.Errorf("ID = %d, want 4", got)
	}
}

func TestCompleteSetsDoneAndTimestamp(t *testing.T) {
	t.Parallel()

	l := newList("a", "b")
	now := time.Date(2026, time.October, 2, 18, 30, 0, 0, time.UTC)
	if err := l.Complete(now, 1); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	got, err := l.Get(1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !got.Done || !got.CompletedAt.Equal(now) {
		t.Errorf("got %+v, want done at %v", got, now)
	}
	if second, err := l.Get(2); err != nil || second.Done {
		t.Errorf("todo 2 = %+v (err %v), want pending", second, err)
	}
}

func TestCompleteUnknownIDIsNoop(t *testing.T) {
	t.Parallel()

	l := newList("a")
	err := l.Complete(time.Now(), 1, 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if l.Todos[0].Done {
		t.Error("todo 1 marked done despite failed batch")
	}
}

func TestReopenClearsTimestamp(t *testing.T) {
	t.Parallel()

	l := newList("a")
	if err := l.Complete(time.Now(), 1); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if err := l.Reopen(1); err != nil {
		t.Fatalf("Reopen: %v", err)
	}
	if got := l.Todos[0]; got.Done || !got.CompletedAt.IsZero() {
		t.Errorf("got %+v, want pending with zero CompletedAt", got)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	l := newList("a", "b", "c")
	n, err := l.Remove(1, 3)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if n != 2 {
		t.Errorf("removed = %d, want 2", n)
	}
	if got, want := ids(l.Todos), []int{2}; !slices.Equal(got, want) {
		t.Errorf("ids = %v, want %v", got, want)
	}
}

func TestRemoveUnknownIDKeepsListIntact(t *testing.T) {
	t.Parallel()

	l := newList("a", "b")
	n, err := l.Remove(1, 42)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if n != 0 {
		t.Errorf("removed = %d, want 0", n)
	}
	if len(l.Todos) != 2 {
		t.Errorf("len = %d, want 2 (list must be unchanged)", len(l.Todos))
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	l := newList("a", "b", "c")
	if err := l.Complete(time.Now(), 2); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	if got, want := ids(l.Filter(false)), []int{1, 3}; !slices.Equal(got, want) {
		t.Errorf("pending = %v, want %v", got, want)
	}
	if got, want := ids(l.Filter(true)), []int{2}; !slices.Equal(got, want) {
		t.Errorf("done = %v, want %v", got, want)
	}
}

func TestClear(t *testing.T) {
	t.Parallel()

	l := newList("a", "b")
	if n := l.Clear(); n != 2 {
		t.Errorf("cleared = %d, want 2", n)
	}
	if len(l.Todos) != 0 {
		t.Errorf("len = %d, want 0", len(l.Todos))
	}
}

func TestParseIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    []int
		wantErr bool
	}{
		{name: "single", args: []string{"3"}, want: []int{3}},
		{name: "comma separated", args: []string{"1,2,3"}, want: []int{1, 2, 3}},
		{name: "multiple args", args: []string{"1,2", " 4 "}, want: []int{1, 2, 4}},
		{name: "empty", args: nil, wantErr: true},
		{name: "not a number", args: []string{"1,x"}, wantErr: true},
		{name: "zero", args: []string{"0"}, wantErr: true},
		{name: "negative", args: []string{"-2"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseIDs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseIDs(%q) = %v, want error", tt.args, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseIDs(%q): %v", tt.args, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ParseIDs(%q) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
