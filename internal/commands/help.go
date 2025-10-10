package commands

import (
	"fmt"
	"sort"
	"strings"
)

type helpCommand struct {
	registry *Registry
}

// NewHelpCommand returns a command that renders inline documentation for other
// commands. The registry is used to enumerate available commands.
func NewHelpCommand(registry *Registry) Command {
	return &helpCommand{registry: registry}
}

func (h *helpCommand) Name() string { return "help" }

func (h *helpCommand) Aliases() []string { return []string{"?"} }

func (h *helpCommand) Description() string {
	return "Show a list of available commands"
}

func (h *helpCommand) Usage() string {
	return "help [command]"
}

func (h *helpCommand) Handle(req Request) (Response, error) {
	if len(req.Args) == 0 {
		lines := []string{"*Available commands:*"}

		cmds := h.registry.Commands()
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name() < cmds[j].Name()
		})

		for _, cmd := range cmds {
			lines = append(lines, fmt.Sprintf("• *%s*: %s", cmd.Name(), cmd.Description()))
		}

		return Response{Text: strings.Join(lines, "\n"), UseMarkdown: true}, nil
	}

	name := req.Args[0]
	if cmd, ok := h.registry.Find(name); ok {
		text := fmt.Sprintf("*%s*\n%s\n_Usage:_ `%s`", cmd.Name(), cmd.Description(), cmd.Usage())
		return Response{Text: text, UseMarkdown: true}, nil
	}

	return Response{Text: fmt.Sprintf("I don't know anything about `%s`", name), UseMarkdown: true}, nil
}
