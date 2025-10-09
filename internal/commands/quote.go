package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type quoteCommand struct{}

const quotePrefix = "quote:"

// NewQuoteCommand returns a command that manages memorable quotes.
func NewQuoteCommand() Command { return &quoteCommand{} }

func (q *quoteCommand) Name() string { return "quote" }

func (q *quoteCommand) Aliases() []string { return []string{"quotes"} }

func (q *quoteCommand) Description() string {
	return "Share or store memorable moments"
}

func (q *quoteCommand) Usage() string {
	return "quote [add <text>]"
}

func (q *quoteCommand) Handle(req Request) (Response, error) {
	if len(req.Args) >= 1 && strings.EqualFold(req.Args[0], "add") {
		if len(req.Args) < 2 {
			return Response{Text: "Usage: quote add <text>"}, nil
		}

		text := strings.Join(req.Args[1:], " ")
		key := fmt.Sprintf("%s%s", quotePrefix, uuid.NewString())
		if err := req.Brain.Set(req.Context, key, fmt.Sprintf("%s — <@%s> (%s)", text, req.User, time.Now().Format(time.RFC822))); err != nil {
			return Response{}, fmt.Errorf("store quote: %w", err)
		}

		return Response{Text: "Saved the moment for posterity."}, nil
	}

	if _, value, err := req.Brain.Random(req.Context, quotePrefix); err == nil {
		return Response{Text: value, UseMarkdown: true}, nil
	}

	return Response{Text: "I don't have any quotes yet. Add one with `quote add <text>`.", UseMarkdown: true}, nil
}
