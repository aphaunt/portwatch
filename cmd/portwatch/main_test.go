package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestMain_MissingConfig verifies the binary exits non-zero when the config
// file does not exist. This is an integration smoke-test that compiles and
// runs the actual binary.
func TestMain_MissingConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping binary integration test in short mode")
	}

	// Build the binary into a temp file.
	tmp, err := os.CreateTemp("", "portwatch-test-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	build := exec.Command("go", "build", "-o", tmp.Name(), ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(tmp.Name(), "-config", "/nonexistent/portwatch.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected non-zero exit, got nil")
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		t.Fatal("binary did not exit within timeout")
	}
}
