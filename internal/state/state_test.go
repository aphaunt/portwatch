package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/state"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	store := state.NewStore(tempStorePath(t))

	snap := state.Snapshot{
		Ports:      []int{80, 443, 8080},
		RecordedAt: time.Now().Truncate(time.Second),
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Ports) != len(snap.Ports) {
		t.Errorf("ports len: got %d, want %d", len(loaded.Ports), len(snap.Ports))
	}
	for i, p := range snap.Ports {
		if loaded.Ports[i] != p {
			t.Errorf("port[%d]: got %d, want %d", i, loaded.Ports[i], p)
		}
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	store := state.NewStore(tempStorePath(t))

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty snapshot, got %v", snap.Ports)
	}
}

func TestStore_Clear(t *testing.T) {
	path := tempStorePath(t)
	store := state.NewStore(path)

	if err := store.Save(state.Snapshot{Ports: []int{22}}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := store.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestStore_ClearIdempotent(t *testing.T) {
	store := state.NewStore(tempStorePath(t))
	if err := store.Clear(); err != nil {
		t.Errorf("Clear on non-existent file should not error: %v", err)
	}
}

func TestStore_LoadInvalidJSON(t *testing.T) {
	path := tempStorePath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	store := state.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
