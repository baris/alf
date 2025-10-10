package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/baris/alf/internal/brain"
	"github.com/baris/alf/internal/commands"
	"github.com/baris/alf/internal/config"
)

// App provides a lightweight, dependency-free runtime for alf.
type App struct {
	cfg      config.Config
	brain    brain.Store
	registry *commands.Registry

	in     io.Reader
	out    io.Writer
	greet  string
	user   string
	prompt string
	once   sync.Once
}

// New constructs an App using the provided configuration.
func New(cfg config.Config) (*App, error) {
	store, err := brain.NewFileStore(cfg.BrainPath)
	if err != nil {
		return nil, fmt.Errorf("init brain: %w", err)
	}

	registry := commands.NewRegistry()
	registry.Register(commands.NewHelpCommand(registry))
	registry.Register(commands.NewPingCommand())
	registry.Register(commands.NewRememberCommand())
	registry.Register(commands.NewRecallCommand())
	registry.Register(commands.NewForgetCommand())
	registry.Register(commands.NewQuoteCommand())
	registry.Register(commands.NewAboutCommand())

	app := &App{
		cfg:      cfg,
		brain:    store,
		registry: registry,
		in:       os.Stdin,
		out:      os.Stdout,
		greet:    cfg.Greeting,
		user:     cfg.DefaultUser,
		prompt:   cfg.Prompt,
	}

	if app.prompt == "" {
		app.prompt = "> "
	}

	if app.greet == "" {
		app.greet = "alf ready. type 'help' to explore commands or 'exit' to quit."
	}

	if app.user == "" {
		app.user = "operator"
	}

	return app, nil
}

// WithIO replaces the default input and output streams. Primarily used for tests.
func (a *App) WithIO(in io.Reader, out io.Writer) *App {
	a.in = in
	a.out = out
	return a
}

// Run starts a simple REPL that executes commands until the context is cancelled
// or the user types "exit" or "quit".
func (a *App) Run(ctx context.Context) error {
	a.once.Do(func() {
		fmt.Fprintln(a.out, a.greet)
	})

	scanner := bufio.NewScanner(a.in)
	for {
		if _, err := fmt.Fprint(a.out, a.prompt); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}
			return nil
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if exitRequested(line) {
			fmt.Fprintln(a.out, "goodbye! 👋")
			return nil
		}

		resp, found, err := a.executeCommand(ctx, line)
		if err != nil {
			fmt.Fprintf(a.out, "error: %v\n", err)
			continue
		}

		if !found {
			fmt.Fprintln(a.out, "I'm still learning. Try 'help' for what I can do.")
			continue
		}

		if resp.Text != "" {
			fmt.Fprintln(a.out, resp.Text)
		}
	}
}

func (a *App) executeCommand(ctx context.Context, text string) (commands.Response, bool, error) {
	name, args := parseCommand(text)
	if name == "" {
		return commands.Response{}, false, nil
	}

	req := commands.Request{
		Context: ctx,
		Brain:   a.brain,
		User:    a.user,
		Text:    text,
		Args:    args,
	}

	cmd, ok := a.registry.Find(name)
	if !ok {
		return commands.Response{}, false, nil
	}

	resp, err := cmd.Handle(req)
	if err != nil {
		return commands.Response{}, true, err
	}

	return resp, true, nil
}

func parseCommand(text string) (string, []string) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return "", nil
	}

	name := strings.ToLower(fields[0])
	args := []string(nil)
	if len(fields) > 1 {
		args = fields[1:]
	}

	return name, args
}

func exitRequested(text string) bool {
	lower := strings.ToLower(text)
	if lower == "exit" || lower == "quit" {
		return true
	}

	return strings.HasPrefix(lower, ":q")
}
