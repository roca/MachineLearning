package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Store persists a List as a JSON document.
type Store struct {
	path string
}

// NewStore returns a Store backed by path.
func NewStore(path string) *Store { return &Store{path: path} }

// Path returns the file backing the store.
func (s *Store) Path() string { return s.path }

// DefaultPath returns the data file location, honouring $TODO_FILE.
func DefaultPath() (string, error) {
	if p := os.Getenv("TODO_FILE"); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate user config dir: %w", err)
	}
	return filepath.Join(dir, "todo-cli", "todos.json"), nil
}

// Load reads the list. A missing file yields an empty list.
func (s *Store) Load() (*List, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return &List{NextID: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", s.path, err)
	}
	l := &List{}
	if err := json.Unmarshal(data, l); err != nil {
		return nil, fmt.Errorf("parse %s: %w", s.path, err)
	}
	if err := l.Validate(); err != nil {
		return nil, fmt.Errorf("validate %s: %w", s.path, err)
	}
	l.NextID = max(l.NextID, l.maxID()+1)
	return l, nil
}

// maxID returns the highest ID in use, or 0 when the list is empty.
func (l *List) maxID() int {
	id := 0
	for _, t := range l.Todos {
		id = max(id, t.ID)
	}
	return id
}

// Save writes the list atomically, creating parent directories as needed.
func (s *Store) Save(l *List) error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return fmt.Errorf("encode todos: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".todos-*.json")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write %s: %w", tmp.Name(), err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", tmp.Name(), err)
	}
	if err := os.Chmod(tmp.Name(), 0o600); err != nil {
		return fmt.Errorf("chmod %s: %w", tmp.Name(), err)
	}
	if err := os.Rename(tmp.Name(), s.path); err != nil {
		return fmt.Errorf("rename %s: %w", s.path, err)
	}
	return nil
}
