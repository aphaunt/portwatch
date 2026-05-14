// Package scanner provides functionality for scanning TCP/UDP ports on a target host.
//
// Basic usage:
//
//	s := scanner.NewScanner("127.0.0.1", 500*time.Millisecond)
//	results, err := s.Scan("tcp", 1, 1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, r := range results {
//		if r.Open {
//			fmt.Printf("Port %d/%s is OPEN\n", r.Port, r.Protocol)
//		}
//	}
//
// The scanner performs a simple dial-based probe: if a connection can be
// established within the configured timeout the port is considered open.
// Each PortState result records the port number, protocol, open/closed status,
// and the timestamp at which it was probed.
package scanner
