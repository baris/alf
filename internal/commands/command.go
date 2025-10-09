package commands

import (
	"context"
	"strings"

	"github.com/baris/alf/internal/brain"
)

// Request wraps the payload required to execute a command.
type Request struct {
	Context context.Context
	Brain   brain.Store
	User    string
	Text    string
	Args    []string
}

// Response represents a command result ready for display to the user.
type Response struct {
	Text        string
	UseMarkdown bool
}

// Command represents a unit of work that can respond to user input.
type Command interface {
	Name() string
	Aliases() []string
	Description() string
	Usage() string
	Handle(req Request) (Response, error)
}

// Registry manages the available commands.
type Registry struct {
	commands map[string]Command
	order    []Command
}

// NewRegistry constructs a registry with no commands.
func NewRegistry() *Registry {
	return &Registry{
		commands: map[string]Command{},
		order:    []Command{},
	}
}

// Register adds a command to the registry. The last registered command wins
// when aliases collide.
func (r *Registry) Register(cmd Command) {
	if cmd == nil {
		return
	}

	r.order = append(r.order, cmd)
	r.commands[strings.ToLower(cmd.Name())] = cmd
	for _, alias := range cmd.Aliases() {
		r.commands[strings.ToLower(alias)] = cmd
	}
}

// Find returns a command for the provided name.
func (r *Registry) Find(name string) (Command, bool) {
	cmd, ok := r.commands[strings.ToLower(name)]
	return cmd, ok
}

// Commands returns the registered commands in the order they were added.
func (r *Registry) Commands() []Command {
	return append([]Command(nil), r.order...)
}
