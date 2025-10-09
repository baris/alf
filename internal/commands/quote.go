package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
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
		key := fmt.Sprintf("%s%s", quotePrefix, randomID())
		value := fmt.Sprintf("%s — %s (%s)", text, req.User, time.Now().Format(time.RFC822))
		if err := req.Brain.Set(req.Context, key, value); err != nil {
			return Response{}, fmt.Errorf("store quote: %w", err)
		}

		return Response{Text: "Saved the moment for posterity."}, nil
	}

	if _, value, err := req.Brain.Random(req.Context, quotePrefix); err == nil {
		return Response{Text: value, UseMarkdown: true}, nil
	}

	return Response{Text: "I don't have any quotes yet. Add one with `quote add <text>`.", UseMarkdown: true}, nil
}

func randomID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
