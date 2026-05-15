package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends alert notifications to a Slack incoming webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

// slackPayload is the JSON body sent to Slack.
type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier that posts to the given Slack
// incoming webhook URL. A default HTTP client with a 10-second timeout is used.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify sends the alert to Slack as a formatted message.
func (s *SlackNotifier) Notify(a Alert) error {
	payload := slackPayload{
		Text: fmt.Sprintf(":rotating_light: *portwatch alert* — %s", a.String()),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}

	return nil
}
