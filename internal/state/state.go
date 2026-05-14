package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Ports     []int     `json:"ports"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Store persists and retrieves port snapshots to/from disk.
type Store struct {
	mu   sync.RWMutex
	path string
}

// NewStore creates a Store backed by the file at path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk, overwriting any existing state.
func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Load reads the most recent snapshot from disk.
// If the file does not exist, an empty Snapshot and nil error are returned.
func (s *Store) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Clear removes the state file from disk.
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
