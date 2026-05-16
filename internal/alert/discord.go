package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DiscordNotifier sends alert notifications to a Discord channel via webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// discordPayload represents the JSON body sent to the Discord webhook API.
type discordPayload struct {
	Content string `json:"content"`
}

// NewDiscordNotifier creates a DiscordNotifier that posts messages to the
// given Discord webhook URL.
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify sends the alert to the configured Discord webhook.
func (d *DiscordNotifier) Notify(a Alert) error {
	payload := discordPayload{
		Content: fmt.Sprintf("**portwatch alert** — %s", a.String()),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: post webhook: %w", err)
	}
	defer resp.Body.Close()

	// Discord returns 204 No Content on success.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}

	return nil
}
