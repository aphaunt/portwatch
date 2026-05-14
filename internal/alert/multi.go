package alert

import "errors"

// MultiNotifier fans out a single Alert to multiple Notifier implementations.
// All notifiers are called even if one returns an error; errors are joined.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that forwards to each provided Notifier.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Add appends a Notifier to the fan-out list.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Notify delivers the alert to every registered Notifier.
// It collects all errors and returns them joined.
func (m *MultiNotifier) Notify(a Alert) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(a); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
