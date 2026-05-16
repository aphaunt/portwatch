// Package alert defines the Notifier interface and concrete implementations
// used by portwatch to deliver port-change alerts through multiple channels.
//
// Built-in notifiers
//
//   - LogNotifier    – writes to a standard Go [log.Logger]
//   - WebhookNotifier – HTTP POST to an arbitrary URL
//   - SlackNotifier   – Slack Incoming Webhook
//   - DiscordNotifier – Discord Webhook
//   - TeamsNotifier   – Microsoft Teams Incoming Webhook
//   - EmailNotifier   – SMTP email delivery
//   - PagerDutyNotifier – PagerDuty Events API v2
//   - OpsGenieNotifier  – OpsGenie Alerts API
//   - SyslogNotifier    – local syslog daemon (LOG_WARNING|LOG_DAEMON)
//
// Notifiers can be composed with [NewMultiNotifier] to fan-out a single
// alert to several backends simultaneously.
package alert
