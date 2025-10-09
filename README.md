# alf (modern edition)

Alf is a modernized Slack bot written in Go. It runs in Socket Mode, keeps a
simple persistent brain for quotes and facts, and exposes a small set of
commands to interact with your team.

## Features

- Socket Mode connection using [`slack-go`](https://github.com/slack-go/slack)
- Structured logging with Uber's [`zap`](https://github.com/uber-go/zap)
- Extensible command registry with built-in commands (`help`, `ping`, `remember`,
  `recall`, `forget`, `quote`, `about`)
- Persistent JSON “brain” for storing facts and quotes

## Getting started

1. Create a Slack application with Socket Mode enabled and grant the bot scope
   permissions your team requires (at minimum `app_mentions:read`, `chat:write`,
   `commands`, and `im:history`).
2. Generate a **Bot Token** (`xoxb-…`) and an **App Level Token** (`xapp-…`).
3. Copy `alf.conf.example` to a new YAML file (e.g. `config.yaml`) and supply the
   token values:

   ```yaml
   app_token: xapp-123
   bot_token: xoxb-456
   brain_path: data/brain.json
   log_level: info
   ```

   The same values can be provided via environment variables `ALF_APP_TOKEN`,
   `ALF_BOT_TOKEN`, `ALF_BRAIN_PATH`, and `ALF_LOG_LEVEL`.
4. Build and run the bot:

   ```bash
   make build
   ./bin/alf --config config.yaml
   ```

   Press `Ctrl+C` to stop the bot.

## Development

This project uses Go modules.

```bash
# Run unit tests
make test

# Format the code
make fmt
```

To add new commands, create a type that implements the `Command` interface in
`internal/commands` and register it in `internal/app/app.go`.
