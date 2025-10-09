package commands

import "runtime"

type aboutCommand struct{}

// NewAboutCommand exposes build information about the bot.
func NewAboutCommand() Command { return &aboutCommand{} }

func (a *aboutCommand) Name() string { return "about" }

func (a *aboutCommand) Aliases() []string { return []string{"version"} }

func (a *aboutCommand) Description() string {
	return "Learn more about this version of alf"
}

func (a *aboutCommand) Usage() string { return "about" }

func (a *aboutCommand) Handle(req Request) (Response, error) {
	text := "I'm the freshly rebooted *alf* powered by Go " + runtime.Version()
	return Response{Text: text, UseMarkdown: true}, nil
}
