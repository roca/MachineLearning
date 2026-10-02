// Package todo implements a todo list with JSON-backed persistence.
package todo

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ErrNotFound is reported when an operation references an unknown todo ID.
var ErrNotFound = errors.New("todo not found")

// Todo is a single list item. ParentID points at the containing todo, or is
// zero for a top-level item.
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Done        bool      `json:"done"`
	ParentID    int       `json:"parent_id,omitzero"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitzero"`
}

// List is a forest of todos with monotonic IDs.
type List struct {
	NextID int    `json:"next_id"`
	Todos  []Todo `json:"todos"`
}

// Node is a todo with its position in the flattened tree.
type Node struct {
	Todo
	Depth int    // nesting level, 0 for top-level items
	Last  bool   // no sibling follows at this level
	Trail []bool // per ancestor level: true draws a vertical bar
}

// Add appends a todo under parentID, or at the top level when parentID is 0.
func (l *List) Add(title string, parentID int, now time.Time) (Todo, error) {
	if parentID != 0 {
		if _, err := l.Get(parentID); err != nil {
			return Todo{}, fmt.Errorf("parent: %w", err)
		}
	}
	t := Todo{ID: l.NextID, Title: title, ParentID: parentID, CreatedAt: now}
	l.NextID++
	l.Todos = append(l.Todos, t)
	return t, nil
}

// Get returns a pointer to the todo with the given ID.
func (l *List) Get(id int) (*Todo, error) {
	i := slices.IndexFunc(l.Todos, func(t Todo) bool { return t.ID == id })
	if i < 0 {
		return nil, fmt.Errorf("id %d: %w", id, ErrNotFound)
	}
	return &l.Todos[i], nil
}

// Roots returns top-level todos sorted by ID.
func (l *List) Roots() []Todo { return l.Children(0) }

// Children returns the direct children of parentID, sorted by ID. ParentID 0
// yields the top-level todos.
func (l *List) Children(parentID int) []Todo {
	kids := make([]Todo, 0, len(l.Todos))
	for _, t := range l.Todos {
		if t.ParentID == parentID {
			kids = append(kids, t)
		}
	}
	slices.SortFunc(kids, func(a, b Todo) int { return cmp.Compare(a.ID, b.ID) })
	return kids
}

// Complete marks every given todo as done. It is atomic: if any ID is
// unknown, no todo is modified.
func (l *List) Complete(now time.Time, ids ...int) error {
	targets, err := l.resolve(ids)
	if err != nil {
		return err
	}
	for _, t := range targets {
		t.Done = true
		t.CompletedAt = now
	}
	return nil
}

// Reopen clears the done flag on every given todo. It is atomic: if any ID
// is unknown, no todo is modified.
func (l *List) Reopen(ids ...int) error {
	targets, err := l.resolve(ids)
	if err != nil {
		return err
	}
	for _, t := range targets {
		t.Done = false
		t.CompletedAt = time.Time{}
	}
	return nil
}

// Remove deletes the given todos together with all of their descendants, so
// no todo is left pointing at a deleted parent. It is atomic: if any ID is
// unknown, nothing is deleted.
func (l *List) Remove(ids ...int) (int, error) {
	if _, err := l.resolve(ids); err != nil {
		return 0, err
	}

	doomed := make([]int, 0, len(ids))
	seen := make(map[int]bool, len(ids))
	var expand func(id int)
	expand = func(id int) {
		if seen[id] {
			return
		}
		seen[id] = true
		doomed = append(doomed, id)
		for _, c := range l.Children(id) {
			expand(c.ID)
		}
	}
	for _, id := range ids {
		expand(id)
	}
	return l.deleteIDs(doomed), nil
}

// Prune deletes the given todos but keeps their descendants, re-parenting
// each child to the deleted todo's own parent. It is atomic: if any ID is
// unknown, nothing is deleted.
func (l *List) Prune(ids ...int) (int, error) {
	targets, err := l.resolve(ids)
	if err != nil {
		return 0, err
	}

	inheritors := make(map[int]int, len(targets))
	for _, t := range targets {
		inheritors[t.ID] = t.ParentID
	}

	// A child can inherit a parent that is itself being pruned, so follow the
	// chain and land on a todo that survives. The step count bounds the walk.
	heir := func(parentID int) int {
		for range inheritors {
			next, pruned := inheritors[parentID]
			if !pruned {
				return parentID
			}
			parentID = next
		}
		return parentID
	}
	for i := range l.Todos {
		if _, pruned := inheritors[l.Todos[i].ID]; pruned {
			continue
		}
		l.Todos[i].ParentID = heir(l.Todos[i].ParentID)
	}
	return l.deleteIDs(ids), nil
}

// Filter returns the todos with the given done state, sorted by ID.
func (l *List) Filter(done bool) []Todo {
	out := slices.DeleteFunc(slices.Clone(l.Todos), func(t Todo) bool { return t.Done != done })
	slices.SortFunc(out, func(a, b Todo) int { return cmp.Compare(a.ID, b.ID) })
	return out
}

// Clear removes all todos and reports how many were removed.
func (l *List) Clear() int {
	n := len(l.Todos)
	l.Todos = nil
	return n
}

// Tree flattens the forest depth-first, siblings sorted by ID. includeDone
// controls whether completed todos are listed; a todo is always kept when one
// of its descendants passes the filter, so filtering never hides a subtree.
func (l *List) Tree(includeDone bool) []Node {
	var nodes []Node
	seen := make(map[int]bool, len(l.Todos))

	var walk func(parentID, depth int, trail []bool)
	walk = func(parentID, depth int, trail []bool) {
		kids := l.Children(parentID)
		for i, kid := range kids {
			if seen[kid.ID] {
				continue // corrupt data: stop a parent cycle from looping
			}
			seen[kid.ID] = true
			last := i == len(kids)-1
			if includeDone || !kid.Done || l.hasVisibleChild(kid.ID, includeDone, seen) {
				nodes = append(nodes, Node{Todo: kid, Depth: depth, Last: last, Trail: trail})
			}
			walk(kid.ID, depth+1, slices.Concat(trail, []bool{!last}))
		}
	}
	walk(0, 0, nil)
	return nodes
}

// hasVisibleChild reports whether parentID has a descendant that passes the
// filter, walking only into subtrees not yet visited.
func (l *List) hasVisibleChild(parentID int, includeDone bool, seen map[int]bool) bool {
	for _, c := range l.Children(parentID) {
		if seen[c.ID] {
			continue
		}
		if includeDone || !c.Done {
			return true
		}
		seen[c.ID] = true
		if l.hasVisibleChild(c.ID, includeDone, seen) {
			return true
		}
	}
	return false
}

// Validate reports the first structural problem in the list: a duplicate or
// invalid ID, an unknown parent, or a parent cycle.
func (l *List) Validate() error {
	byID := make(map[int]Todo, len(l.Todos))
	for _, t := range l.Todos {
		switch {
		case t.ID < 1:
			return fmt.Errorf("todo %d: id must be positive", t.ID)
		case t.ParentID == t.ID:
			return fmt.Errorf("todo %d: cannot be its own parent", t.ID)
		}
		if _, dup := byID[t.ID]; dup {
			return fmt.Errorf("todo %d: duplicate id", t.ID)
		}
		byID[t.ID] = t
	}

	for _, t := range l.Todos {
		if t.ParentID == 0 {
			continue
		}
		if _, ok := byID[t.ParentID]; !ok {
			return fmt.Errorf("todo %d: parent %d not found", t.ID, t.ParentID)
		}
		for id, steps := t.ParentID, 0; id != 0; id, steps = byID[id].ParentID, steps+1 {
			if id == t.ID || steps > len(byID) {
				return fmt.Errorf("todo %d: parent cycle", t.ID)
			}
		}
	}
	return nil
}

// ParseIDs parses comma-separated positive integers such as "1,2, 3".
func ParseIDs(args []string) ([]int, error) {
	var ids []int
	for arg := range strings.SplitSeq(strings.Join(args, ","), ",") {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		id, err := strconv.Atoi(arg)
		if err != nil || id < 1 {
			return nil, fmt.Errorf("invalid id %q", arg)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, errors.New("no ids given")
	}
	return ids, nil
}

// resolve looks up every ID before any mutation happens.
func (l *List) resolve(ids []int) ([]*Todo, error) {
	targets := make([]*Todo, 0, len(ids))
	for _, id := range ids {
		t, err := l.Get(id)
		if err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (l *List) deleteIDs(ids []int) int {
	kept := slices.DeleteFunc(slices.Clone(l.Todos), func(t Todo) bool {
		return slices.Contains(ids, t.ID)
	})
	removed := len(l.Todos) - len(kept)
	l.Todos = kept
	return removed
}
