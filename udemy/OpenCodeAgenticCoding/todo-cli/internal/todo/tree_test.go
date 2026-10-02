package todo

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

// sample builds the tree used by several tests:
//
//	1 alpha
//	├── 2 bravo
//	│   ├── 4 delta
//	│   └── 5 echo
//	└── 3 charlie
//	6 foxtrot
func sample(t *testing.T) *List {
	t.Helper()

	l := newList("alpha")   // 1
	add(t, l, "bravo", 1)   // 2
	add(t, l, "charlie", 1) // 3
	add(t, l, "delta", 2)   // 4
	add(t, l, "echo", 2)    // 5
	add(t, l, "foxtrot", 0) // 6
	return l
}

func TestAddWithParent(t *testing.T) {
	t.Parallel()

	l := newList("alpha")
	child := add(t, l, "bravo", 1)
	if child.ParentID != 1 {
		t.Errorf("ParentID = %d, want 1", child.ParentID)
	}
	if got := add(t, l, "solo", 0); got.ParentID != 0 {
		t.Errorf("ParentID = %d, want 0", got.ParentID)
	}
}

func TestAddWithUnknownParent(t *testing.T) {
	t.Parallel()

	l := newList("alpha")
	_, err := l.Add("orphan", 99, created)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if len(l.Todos) != 1 {
		t.Errorf("len = %d, want 1 (nothing added)", len(l.Todos))
	}
}

func TestAddCannotNestUnderChild(t *testing.T) {
	t.Parallel()

	// A parent must already exist, so a brand new todo can never be its own
	// ancestor; this guards the invariant that Tree relies on.
	l := newList("alpha")
	child := add(t, l, "bravo", 1)
	if _, err := l.Add("charlie", child.ID, created); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := l.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestChildrenAndRoots(t *testing.T) {
	t.Parallel()

	l := sample(t)
	if got, want := ids(l.Roots()), []int{1, 6}; !slices.Equal(got, want) {
		t.Errorf("Roots = %v, want %v", got, want)
	}
	if got, want := ids(l.Children(1)), []int{2, 3}; !slices.Equal(got, want) {
		t.Errorf("Children(1) = %v, want %v", got, want)
	}
	if got, want := ids(l.Children(2)), []int{4, 5}; !slices.Equal(got, want) {
		t.Errorf("Children(2) = %v, want %v", got, want)
	}
	if got := l.Children(3); len(got) != 0 {
		t.Errorf("Children(3) = %v, want none", got)
	}
}

func TestTreeOrderAndShape(t *testing.T) {
	t.Parallel()

	type shape struct {
		id    int
		depth int
		last  bool
		trail []bool
	}
	want := []shape{
		{id: 1, depth: 0, last: false, trail: nil},
		{id: 2, depth: 1, last: false, trail: []bool{true}},
		{id: 4, depth: 2, last: false, trail: []bool{true, true}},
		{id: 5, depth: 2, last: true, trail: []bool{true, true}},
		{id: 3, depth: 1, last: true, trail: []bool{true}},
		{id: 6, depth: 0, last: true, trail: nil},
	}

	got := sample(t).Tree(true)
	if len(got) != len(want) {
		t.Fatalf("Tree returned %d nodes, want %d: %+v", len(got), len(want), got)
	}
	for i, n := range got {
		if n.ID != want[i].id || n.Depth != want[i].depth || n.Last != want[i].last || !slices.Equal(n.Trail, want[i].trail) {
			t.Errorf("node %d = %+v (trail %v), want %+v", i, n, n.Trail, want[i])
		}
	}
}

func TestTreeTrailsAreNotAliased(t *testing.T) {
	t.Parallel()

	// Siblings legitimately share their parent's trail, but two nodes with
	// different trails must never share a backing array: a walk that appended
	// in place would rewrite an earlier node's prefix.
	nodes := sample(t).Tree(true)
	for i, a := range nodes {
		for j, b := range nodes {
			if i == j || len(a.Trail) == 0 || len(b.Trail) == 0 || slices.Equal(a.Trail, b.Trail) {
				continue
			}
			if &a.Trail[0] == &b.Trail[0] {
				t.Errorf("todos %d and %d share a trail array with different contents %v / %v",
					a.ID, b.ID, a.Trail, b.Trail)
			}
		}
	}
}

func TestTreeFiltersDoneButKeepsContext(t *testing.T) {
	t.Parallel()

	l := sample(t)
	if err := l.Complete(created, 1, 4); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	got := idsOfNodes(l.Tree(false))
	want := []int{1, 2, 5, 3, 6} // 4 is done; 1 stays as context for 2 and 3
	if !slices.Equal(got, want) {
		t.Errorf("Tree(false) = %v, want %v", got, want)
	}

	if err := l.Complete(created, 2, 5); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	got = idsOfNodes(l.Tree(false))
	want = []int{1, 3, 6} // subtree 2 is fully done, so it disappears
	if !slices.Equal(got, want) {
		t.Errorf("Tree(false) = %v, want %v", got, want)
	}
}

func TestRemoveCascadesToDescendants(t *testing.T) {
	t.Parallel()

	l := sample(t)
	n, err := l.Remove(2)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if n != 3 {
		t.Errorf("removed = %d, want 3 (#2, #4, #5)", n)
	}
	if got, want := idsOfNodes(l.Tree(true)), []int{1, 3, 6}; !slices.Equal(got, want) {
		t.Errorf("tree = %v, want %v", got, want)
	}
	if err := l.Validate(); err != nil {
		t.Fatalf("Validate after cascade: %v", err)
	}
}

func TestRemoveLeafKeepsSiblings(t *testing.T) {
	t.Parallel()

	l := sample(t)
	n, err := l.Remove(4)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if n != 1 {
		t.Errorf("removed = %d, want 1", n)
	}
	if got, want := idsOfNodes(l.Tree(true)), []int{1, 2, 5, 3, 6}; !slices.Equal(got, want) {
		t.Errorf("tree = %v, want %v", got, want)
	}
}

func TestRemoveCascadesAcrossLevels(t *testing.T) {
	t.Parallel()

	l := sample(t)
	n, err := l.Remove(1)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if n != 5 {
		t.Errorf("removed = %d, want 5", n)
	}
	if got, want := idsOfNodes(l.Tree(true)), []int{6}; !slices.Equal(got, want) {
		t.Errorf("tree = %v, want %v", got, want)
	}
}

func TestPruneReparentsChildren(t *testing.T) {
	t.Parallel()

	l := sample(t)
	n, err := l.Prune(2)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if n != 1 {
		t.Errorf("pruned = %d, want 1", n)
	}
	if got, want := idsOfNodes(l.Tree(true)), []int{1, 3, 4, 5, 6}; !slices.Equal(got, want) {
		t.Errorf("tree = %v, want %v", got, want)
	}
	for _, id := range []int{4, 5} {
		got, err := l.Get(id)
		if err != nil {
			t.Fatalf("Get(%d): %v", id, err)
		}
		if got.ParentID != 1 {
			t.Errorf("todo %d ParentID = %d, want 1 (grandparent)", id, got.ParentID)
		}
	}
	if err := l.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestPruneTopLevelPromotesToTopLevel(t *testing.T) {
	t.Parallel()

	l := sample(t)
	if _, err := l.Prune(1); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if got, want := idsOfNodes(l.Tree(true)), []int{2, 4, 5, 3, 6}; !slices.Equal(got, want) {
		t.Errorf("tree = %v, want %v", got, want)
	}
}

func TestPruneUnknownIDIsNoop(t *testing.T) {
	t.Parallel()

	l := sample(t)
	if _, err := l.Prune(1, 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if len(l.Todos) != 6 {
		t.Errorf("len = %d, want 6 (nothing changed)", len(l.Todos))
	}
	for i, td := range l.Todos {
		if i == 1 && td.ParentID != 1 {
			t.Errorf("todo %d re-parented despite failed Prune", td.ID)
		}
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		list List
		want string
	}{
		{name: "empty", list: List{NextID: 1}},
		{name: "flat", list: List{NextID: 3, Todos: []Todo{{ID: 1}, {ID: 2}}}, want: ""},
		{
			name: "nested",
			list: List{NextID: 3, Todos: []Todo{{ID: 1}, {ID: 2, ParentID: 1}, {ID: 3, ParentID: 2}}},
			want: "",
		},
		{
			name: "self parent",
			list: List{NextID: 2, Todos: []Todo{{ID: 1, ParentID: 1}}},
			want: "own parent",
		},
		{
			name: "missing parent",
			list: List{NextID: 2, Todos: []Todo{{ID: 1, ParentID: 7}}},
			want: "parent 7 not found",
		},
		{
			name: "two cycle",
			list: List{NextID: 3, Todos: []Todo{{ID: 1, ParentID: 2}, {ID: 2, ParentID: 1}}},
			want: "parent cycle",
		},
		{
			name: "three cycle",
			list: List{NextID: 4, Todos: []Todo{{ID: 1, ParentID: 3}, {ID: 2, ParentID: 1}, {ID: 3, ParentID: 2}}},
			want: "parent cycle",
		},
		{
			name: "duplicate id",
			list: List{NextID: 2, Todos: []Todo{{ID: 1}, {ID: 1}}},
			want: "duplicate id",
		},
		{
			name: "zero id",
			list: List{NextID: 1, Todos: []Todo{{ID: 0}}},
			want: "id must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.list.Validate()
			if tt.want == "" {
				if err != nil {
					t.Fatalf("Validate: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Validate = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestTreeSurvivesCorruptCycle(t *testing.T) {
	t.Parallel()

	// A cycle cannot be built through Add, but a hand-edited file can hold one.
	l := &List{NextID: 3, Todos: []Todo{
		{ID: 1, ParentID: 2},
		{ID: 2, ParentID: 1},
	}}
	if err := l.Validate(); err == nil {
		t.Fatal("Validate accepted a parent cycle")
	}

	done := make(chan []Node, 1)
	go func() { done <- l.Tree(true) }()
	select {
	case <-done: // traversal terminated
	case <-time.After(time.Second):
		t.Fatal("Tree did not terminate on a parent cycle")
	}
}

func idsOfNodes(nodes []Node) []int {
	out := make([]int, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, n.ID)
	}
	return out
}
