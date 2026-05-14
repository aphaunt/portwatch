package state

// Diff describes the changes between two port snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// IsEmpty reports whether the diff contains no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare returns a Diff describing ports that were opened or closed
// between the previous and current snapshots.
func Compare(prev, curr Snapshot) Diff {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	var opened, closed []int

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}

	return Diff{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
