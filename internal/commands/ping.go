package commands

import "time"

type pingCommand struct{}

// NewPingCommand returns a command that responds with a simple heartbeat.
func NewPingCommand() Command { return &pingCommand{} }

func (p *pingCommand) Name() string { return "ping" }

func (p *pingCommand) Aliases() []string { return []string{"pong"} }

func (p *pingCommand) Description() string {
	return "Check if alf is awake"
}

func (p *pingCommand) Usage() string { return "ping" }

func (p *pingCommand) Handle(req Request) (Response, error) {
	return Response{Text: "Pong! " + time.Now().Format(time.RFC1123Z)}, nil
}
