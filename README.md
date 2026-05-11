# portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build ./...
```

## Usage

Start the daemon with a configuration file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
  path: /var/log/portwatch.log
whitelist:
  - 22
  - 80
  - 443
```

Once running, portwatch will scan open ports at the specified interval and emit an alert whenever a port outside the whitelist is opened or an expected port disappears.

Example alert output:

```
[ALERT] 2024-01-15T10:42:00Z new port detected: 8080 (PID 3821, nginx)
[ALERT] 2024-01-15T10:43:00Z port closed unexpectedly: 443
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `./config.yaml` | Path to configuration file |
| `--interval` | `30s` | Polling interval |
| `--verbose` | `false` | Enable verbose logging |

## License

MIT © [yourusername](https://github.com/yourusername)