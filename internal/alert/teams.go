package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TeamsNotifier sends alert notifications to a Microsoft Teams channel
// via an incoming webhook URL.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

// teamsPayload represents the Adaptive Card payload sent to Teams.
type teamsPayload struct {
	Type       string        `json:"type"`
	Attachments []teamsAttachment `json:"attachments"`
}

type teamsAttachment struct {
	ContentType string      `json:"contentType"`
	Content     teamsCard   `json:"content"`
}

type teamsCard struct {
	Schema  string       `json:"$schema"`
	Type    string       `json:"type"`
	Version string       `json:"version"`
	Body    []teamsBlock `json:"body"`
}

type teamsBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap"`
}

// NewTeamsNotifier creates a TeamsNotifier that posts to the given
// Microsoft Teams incoming webhook URL.
func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends the alert to the configured Teams webhook.
func (t *TeamsNotifier) Notify(a Alert) error {
	payload := teamsPayload{
		Type: "message",
		Attachments: []teamsAttachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: teamsCard{
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:    "AdaptiveCard",
					Version: "1.4",
					Body: []teamsBlock{
						{Type: "TextBlock", Text: "**portwatch alert**", Wrap: true},
						{Type: "TextBlock", Text: a.String(), Wrap: true},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
