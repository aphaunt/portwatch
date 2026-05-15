package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends alert payloads to an HTTP endpoint via POST.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Port      int    `json:"port"`
	Host      string `json:"host"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
// timeout controls the HTTP client deadline per request.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	return &WebhookNotifier{
		url: url,
		client: &http.Client{Timeout: timeout},
	}
}

// Notify serialises the Alert as JSON and POSTs it to the configured URL.
func (w *WebhookNotifier) Notify(a Alert) error {
	payload := webhookPayload{
		Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
		Event:     a.Event.String(),
		Port:      a.Port,
		Host:      a.Host,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}
	return nil
}
