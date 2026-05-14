// Package monitor builds on top of the scanner package to provide stateful
// port monitoring. It remembers the port states observed in the previous scan
// and surfaces only the changes (ports that opened or closed) between
// consecutive scans.
//
// Typical usage:
//
//	s, _ := scanner.NewScanner("127.0.0.1", 500)
//	m := monitor.New(s)
//
//	for {
//		changes, err := m.Scan(1024, 65535)
//		if err != nil {
//			log.Println("scan error:", err)
//		}
//		for _, c := range changes {
//			log.Println(c)
//		}
//		time.Sleep(30 * time.Second)
//	}
package monitor
