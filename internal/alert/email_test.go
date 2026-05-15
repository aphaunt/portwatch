package alert

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

// startFakeSMTP listens on a random TCP port, accepts one connection, reads the
// conversation into a string, then closes. It returns the address and a channel
// that delivers the received data once the connection is closed.
func startFakeSMTP(t *testing.T) (addr string, received <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startFakeSMTP: listen: %v", err)
	}
	ch := make(chan string, 1)
	go func() {
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			ch <- ""
			return
		}
		defer conn.Close()
		// Respond with minimal SMTP greeting so net/smtp doesn't bail out.
		_, _ = conn.Write([]byte("220 fake SMTP ready\r\n"))
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
		data, _ := io.ReadAll(conn)
		ch <- string(data)
	}()
	return ln.Addr().String(), ch
}

func TestEmailNotifier_NoRecipients(t *testing.T) {
	n := NewEmailNotifier(EmailConfig{
		SMTPHost: "127.0.0.1",
		SMTPPort: 2525,
		From:     "watch@example.com",
		To:       nil,
	})
	err := n.Notify(Alert{Kind: KindOpened, Port: 8080})
	if err == nil {
		t.Fatal("expected error for missing recipients, got nil")
	}
	if !strings.Contains(err.Error(), "no recipients") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEmailNotifier_UnreachableHost(t *testing.T) {
	n := NewEmailNotifier(EmailConfig{
		SMTPHost: "127.0.0.1",
		SMTPPort: 19999, // nothing listening here
		From:     "watch@example.com",
		To:       []string{"ops@example.com"},
	})
	err := n.Notify(Alert{Kind: KindOpened, Port: 443})
	if err == nil {
		t.Fatal("expected error for unreachable SMTP host, got nil")
	}
	if !strings.Contains(err.Error(), "send failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
