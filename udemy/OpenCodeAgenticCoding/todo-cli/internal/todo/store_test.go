package todo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestLoadMissingFileReturnsEmptyList(t *testing.T) {
	t.Parallel()

	l, err := NewStore(filepath.Join(t.TempDir(), "todos.json")).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if l.NextID != 1 || len(l.Todos) != 0 {
		t.Errorf("got %+v, want empty list with NextID 1", l)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "todos.json")
	created := time.Date(2026, time.October, 2, 9, 0, 0, 0, time.UTC)

	want := &List{NextID: 3}
	want.Todos = append(want.Todos,
		Todo{ID: 1, Title: "write docs", CreatedAt: created},
		Todo{ID: 2, Title: "ship v1", Done: true, CreatedAt: created, CompletedAt: created.Add(time.Hour)},
	)

	s := NewStore(path)
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !slices.Equal(got.Todos, want.Todos) || got.NextID != want.NextID {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestSaveOmitsZeroCompletedAt(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	l := &List{NextID: 2}
	if _, err := l.Add("a", 0, time.Now()); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := NewStore(path).Save(l); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	fields := raw["todos"].([]any)[0].(map[string]any)
	if _, ok := fields["completed_at"]; ok {
		t.Error("completed_at present for a pending todo, want omitted")
	}
}

func TestSaveLeavesNoTempFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "todos.json"))
	if err := s.Save(&List{NextID: 1}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "todos.json" {
		t.Errorf("dir contains %v, want only todos.json", entries)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := NewStore(path).Load(); err == nil {
		t.Fatal("Load of corrupt file succeeded, want error")
	}
}

func TestSaveLoadNestedTodos(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	created := time.Date(2026, time.October, 2, 9, 0, 0, 0, time.UTC)

	want := &List{NextID: 4}
	want.Todos = append(want.Todos,
		Todo{ID: 1, Title: "alpha", CreatedAt: created},
		Todo{ID: 2, Title: "bravo", ParentID: 1, CreatedAt: created},
		Todo{ID: 3, Title: "charlie", ParentID: 2, Done: true, CreatedAt: created},
	)

	s := NewStore(path)
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), `"parent_id": 1`) {
		t.Errorf("file missing parent_id:\n%s", data)
	}
	if strings.Contains(string(data), `"parent_id": 0`) {
		t.Errorf("file stores a redundant zero parent_id:\n%s", data)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !slices.Equal(got.Todos, want.Todos) {
		t.Errorf("got %+v, want %+v", got.Todos, want.Todos)
	}
}

func TestLoadRejectsDanglingParent(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	body := `{"next_id":2,"todos":[{"id":1,"title":"a","parent_id":7}]}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := NewStore(path).Load(); err == nil {
		t.Fatal("Load accepted a todo with a missing parent, want error")
	}
}

func TestLoadRejectsParentCycle(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	body := `{"next_id":3,"todos":[{"id":1,"parent_id":2},{"id":2,"parent_id":1}]}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := NewStore(path).Load(); err == nil {
		t.Fatal("Load accepted a parent cycle, want error")
	}
}

func TestLoadRepairsNextID(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "todos.json")
	body := `{"next_id":2,"todos":[{"id":9,"title":"a"}]}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	l, err := NewStore(path).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if l.NextID != 10 {
		t.Errorf("NextID = %d, want 10 (must not reuse #9)", l.NextID)
	}
}

func TestDefaultPathHonoursEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.json")
	t.Setenv("TODO_FILE", path)

	got, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath: %v", err)
	}
	if got != path {
		t.Errorf("DefaultPath = %q, want %q", got, path)
	}
}

func TestGetUnknownID(t *testing.T) {
	t.Parallel()

	l := &List{NextID: 1}
	_, err := l.Get(7)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
