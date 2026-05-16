// Package alert defines the Notifier interface and provides multiple
// notification backends for portwatch.
//
// Supported backends:
//
//	- LogNotifier    – writes alerts to a standard logger
//	- WebhookNotifier – HTTP POST to a generic webhook endpoint
//	- SlackNotifier   – Slack incoming webhook
//	- DiscordNotifier – Discord webhook
//	- TeamsNotifier   – Microsoft Teams incoming webhook (Adaptive Card)
//	- EmailNotifier   – SMTP email
//	- PagerDutyNotifier – PagerDuty Events API v2
//	- OpsGenieNotifier  – OpsGenie Alerts API
//
// Multiple notifiers can be composed with NewMultiNotifier, which fans
// out each alert to every registered backend and collects any errors.
//
// Example:
//
//	multi := alert.NewMultiNotifier()
//	multi.Add(alert.NewLogNotifier(log.Default()))
//	multi.Add(alert.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK")))
//	multi.Add(alert.NewTeamsNotifier(os.Getenv("TEAMS_WEBHOOK")))
package alert
