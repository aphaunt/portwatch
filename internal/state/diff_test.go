package state_test

import (
	"sort"
	"testing"

	"github.com/yourorg/portwatch/internal/state"
)

func sortedInts(s []int) []int {
	out := make([]int, len(s))
	copy(out, s)
	sort.Ints(out)
	return out
}

func TestCompare_Opened(t *testing.T) {
	prev := state.Snapshot{Ports: []int{80}}
	curr := state.Snapshot{Ports: []int{80, 443}}

	diff := state.Compare(prev, curr)

	if got := sortedInts(diff.Opened); len(got) != 1 || got[0] != 443 {
		t.Errorf("Opened: got %v, want [443]", got)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("Closed: got %v, want []", diff.Closed)
	}
}

func TestCompare_Closed(t *testing.T) {
	prev := state.Snapshot{Ports: []int{80, 443}}
	curr := state.Snapshot{Ports: []int{80}}

	diff := state.Compare(prev, curr)

	if len(diff.Opened) != 0 {
		t.Errorf("Opened: got %v, want []", diff.Opened)
	}
	if got := sortedInts(diff.Closed); len(got) != 1 || got[0] != 443 {
		t.Errorf("Closed: got %v, want [443]", got)
	}
}

func TestCompare_NoChange(t *testing.T) {
	snap := state.Snapshot{Ports: []int{22, 80, 443}}
	diff := state.Compare(snap, snap)

	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}

func TestCompare_BothEmpty(t *testing.T) {
	diff := state.Compare(state.Snapshot{}, state.Snapshot{})
	if !diff.IsEmpty() {
		t.Error("expected empty diff for two empty snapshots")
	}
}

func TestCompare_OpenedAndClosed(t *testing.T) {
	prev := state.Snapshot{Ports: []int{22, 80}}
	curr := state.Snapshot{Ports: []int{80, 443}}

	diff := state.Compare(prev, curr)

	if got := sortedInts(diff.Opened); len(got) != 1 || got[0] != 443 {
		t.Errorf("Opened: got %v, want [443]", got)
	}
	if got := sortedInts(diff.Closed); len(got) != 1 || got[0] != 22 {
		t.Errorf("Closed: got %v, want [22]", got)
	}
}
