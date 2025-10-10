# alf (console edition)

This project revives the original alf bot as a dependency-free Go application
that runs entirely in your terminal. It keeps a persistent JSON “brain”, exposes
an extensible command registry, and offers a friendly REPL so you can interact
with alf without any third-party services.

## Features

- Interactive terminal experience with a helpful prompt and greeting
- Persistent JSON brain for remembering facts and quotes
- Built-in commands: `help`, `ping`, `remember`, `recall`, `forget`, `quote`, and
  `about`
- Environment variable and JSON configuration support

## Getting started

1. Build the binary:

   ```bash
   make build
   ```

2. (Optional) Copy `alf.conf.example` to `alf.conf` and tweak the settings. The
   file uses JSON and supports the following fields:

   ```json
   {
     "brain_path": "data/brain.json",
     "default_user": "operator",
     "greeting": "alf ready. type 'help' to explore commands or 'exit' to quit.",
     "prompt": "> "
   }
   ```

   The same values can be supplied via environment variables:

   - `ALF_BRAIN_PATH`
   - `ALF_DEFAULT_USER`
   - `ALF_GREETING`
   - `ALF_PROMPT`

3. Run alf:

   ```bash
   ./bin/alf --config alf.conf
   ```

   Use `--brain` to override the storage path or `--user` to change the name
   recorded when storing memories. Type `exit` (or press `Ctrl+C`) to quit.

## Development

This project uses Go modules.

```bash
# Run unit tests (there currently aren't any, but this verifies the build)
make test

# Format the code
make fmt
```

To add new commands, create a type that implements the `Command` interface in
`internal/commands` and register it in `internal/app/app.go`.
