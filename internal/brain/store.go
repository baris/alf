package brain

import "context"

// Store describes the persistence layer used by alf to remember user-provided
// facts, quotes and other snippets.
type Store interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context, prefix string) ([]string, error)
	Random(ctx context.Context, prefix string) (string, string, error)
}
