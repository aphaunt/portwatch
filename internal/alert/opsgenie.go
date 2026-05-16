package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const opsGenieAlertURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey   string
	client   *http.Client
	alertURL string
}

type ogPayload struct {
	Message  string   `json:"message"`
	Alias    string   `json:"alias"`
	Source   string   `json:"source"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
}

// NewOpsGenieNotifier creates an OpsGenieNotifier with the given API key.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 10 * time.Second},
		alertURL: opsGenieAlertURL,
	}
}

// Notify sends the alert to OpsGenie.
func (o *OpsGenieNotifier) Notify(a Alert) error {
	body := ogPayload{
		Message:  a.String(),
		Alias:    fmt.Sprintf("portwatch-port-%d", a.Port),
		Source:   "portwatch",
		Priority: "P3",
		Tags:     []string{"portwatch", string(a.Kind)},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.alertURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
