// Package main is the entry point for the portwatch daemon.
//
// portwatch is a lightweight daemon that periodically scans configured TCP
// port ranges, compares the results against a persisted baseline, and emits
// alerts whenever ports are unexpectedly opened or closed.
//
// Usage:
//
//	portwatch [-config <path>]
//
// Flags:
//
//	-config  Path to the JSON configuration file (default: portwatch.json).
//
// The daemon runs until it receives SIGINT or SIGTERM, at which point it
// performs a clean shutdown and persists the latest port state to disk so
// that the next invocation can diff against it correctly.
package main
