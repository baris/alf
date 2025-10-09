package commands

import (
	"fmt"
	"strings"
)

type rememberCommand struct{}

type recallCommand struct{}

type forgetCommand struct{}

const memoryPrefix = "memory:"

// NewRememberCommand creates a command that persists key/value pairs supplied by users.
func NewRememberCommand() Command { return &rememberCommand{} }

// NewRecallCommand creates a command that retrieves stored facts.
func NewRecallCommand() Command { return &recallCommand{} }

// NewForgetCommand creates a command that deletes stored facts.
func NewForgetCommand() Command { return &forgetCommand{} }

func (r *rememberCommand) Name() string { return "remember" }

func (r *rememberCommand) Aliases() []string { return []string{"store"} }

func (r *rememberCommand) Description() string {
	return "Remember something for later"
}

func (r *rememberCommand) Usage() string {
	return "remember <key> <value>"
}

func (r *rememberCommand) Handle(req Request) (Response, error) {
	if len(req.Args) < 2 {
		return Response{Text: "Usage: remember <key> <value>"}, nil
	}

	key := strings.ToLower(req.Args[0])
	value := strings.Join(req.Args[1:], " ")

	if err := req.Brain.Set(req.Context, memoryPrefix+key, value); err != nil {
		return Response{}, fmt.Errorf("save memory: %w", err)
	}

	return Response{Text: fmt.Sprintf("Okay, I'll remember `%s`.", key), UseMarkdown: true}, nil
}

func (r *recallCommand) Name() string { return "recall" }

func (r *recallCommand) Aliases() []string { return []string{"remembered", "lookup"} }

func (r *recallCommand) Description() string {
	return "Recall a stored fact"
}

func (r *recallCommand) Usage() string {
	return "recall <key>"
}

func (r *recallCommand) Handle(req Request) (Response, error) {
	if len(req.Args) == 0 {
		return Response{Text: "Usage: recall <key>"}, nil
	}

	key := strings.ToLower(req.Args[0])
	if value, ok, err := req.Brain.Get(req.Context, memoryPrefix+key); err != nil {
		return Response{}, fmt.Errorf("lookup memory: %w", err)
	} else if ok {
		return Response{Text: fmt.Sprintf("%s", value)}, nil
	}

	return Response{Text: fmt.Sprintf("I don't remember `%s`.", key), UseMarkdown: true}, nil
}

func (f *forgetCommand) Name() string { return "forget" }

func (f *forgetCommand) Aliases() []string { return []string{"erase"} }

func (f *forgetCommand) Description() string {
	return "Make alf forget a stored fact"
}

func (f *forgetCommand) Usage() string {
	return "forget <key>"
}

func (f *forgetCommand) Handle(req Request) (Response, error) {
	if len(req.Args) == 0 {
		return Response{Text: "Usage: forget <key>"}, nil
	}

	key := strings.ToLower(req.Args[0])
	if err := req.Brain.Delete(req.Context, memoryPrefix+key); err != nil {
		return Response{}, fmt.Errorf("forget memory: %w", err)
	}

	return Response{Text: fmt.Sprintf("Forgot `%s`.", key), UseMarkdown: true}, nil
}
