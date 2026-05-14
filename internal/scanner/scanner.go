package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	ScannedAt time.Time
}

// Scanner scans a host for open ports within a given range.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// NewScanner creates a new Scanner for the given host.
func NewScanner(host string, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: timeout,
	}
}

// Scan checks all ports in [startPort, endPort] over the given protocol ("tcp" or "udp").
// Returns a slice of PortState for every port in the range.
func (s *Scanner) Scan(protocol string, startPort, endPort int) ([]PortState, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	results := make([]PortState, 0, endPort-startPort+1)

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		open := isPortOpen(protocol, address, s.Timeout)
		results = append(results, PortState{
			Port:      port,
			Protocol:  protocol,
			Open:      open,
			ScannedAt: time.Now(),
		})
	}

	return results, nil
}

// isPortOpen attempts a connection to determine if the port is open.
func isPortOpen(protocol, address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout(protocol, address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
