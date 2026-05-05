# cronwrap

Lightweight wrapper for cron jobs that adds logging, alerting, and run-history tracking.

---

## Installation

```bash
go install github.com/yourusername/cronwrap@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronwrap.git && cd cronwrap && go build ./...
```

---

## Usage

Wrap any command by prefixing it with `cronwrap`:

```bash
cronwrap --job "backup" --alert-on-fail --notify slack ./scripts/backup.sh
```

Example crontab entry:

```cron
0 2 * * * cronwrap --job "nightly-backup" --timeout 30m --log /var/log/cronwrap ./scripts/backup.sh
```

**Common flags:**

| Flag | Description |
|------|-------------|
| `--job` | Human-readable name for the job |
| `--timeout` | Kill the job if it exceeds this duration |
| `--log` | Directory to write structured log files |
| `--alert-on-fail` | Send an alert if the job exits non-zero |
| `--notify` | Notification channel (`slack`, `email`, `webhook`) |
| `--history` | Number of past runs to retain (default: `50`) |

Run history is stored in `~/.cronwrap/history` as JSON and can be reviewed with:

```bash
cronwrap history --job "nightly-backup" --last 10
```

---

## Configuration

Place a `cronwrap.yaml` in your project root or `~/.config/cronwrap/` to set defaults:

```yaml
log_dir: /var/log/cronwrap
history_limit: 100
notify: slack
slack_webhook: https://hooks.slack.com/services/...
```

---

## License

MIT © 2024 Your Name