// Package state provides persistence and diffing of port scan snapshots.
//
// A [Snapshot] captures the set of open ports observed at a moment in time.
// [Store] serialises snapshots to a JSON file on disk so that portwatch can
// detect changes across process restarts.
//
// [Compare] computes the difference between two snapshots, returning the ports
// that were newly opened or closed between observations. This is the primary
// mechanism by which portwatch determines whether an alert should be raised.
//
// Typical usage:
//
//	store := state.NewStore("/var/lib/portwatch/state.json")
//
//	prev, _ := store.Load()
//	curr := state.Snapshot{Ports: scanResult, RecordedAt: time.Now()}
//
//	if diff := state.Compare(prev, curr); !diff.IsEmpty() {
//		// notify about diff.Opened / diff.Closed
//	}
//
//	store.Save(curr)
package state
