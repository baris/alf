package brain

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// fileStore is a threadsafe JSON backed implementation of Store.
type fileStore struct {
	path string
	mu   sync.RWMutex
	data map[string]string
}

// NewFileStore creates a Store that persists data to path.
func NewFileStore(path string) (Store, error) {
	fs := &fileStore{
		path: filepath.Clean(path),
		data: map[string]string{},
	}

	if err := fs.load(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())

	return fs, nil
}

func (s *fileStore) Get(_ context.Context, key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	return value, ok, nil
}

func (s *fileStore) Set(_ context.Context, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return s.persistLocked()
}

func (s *fileStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.persistLocked()
}

func (s *fileStore) Keys(_ context.Context, prefix string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	return keys, nil
}

func (s *fileStore) Random(_ context.Context, prefix string) (string, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}

	if len(keys) == 0 {
		return "", "", errors.New("no entries found")
	}

	idx := rand.Intn(len(keys))
	key := keys[idx]
	return key, s.data[key], nil
}

func (s *fileStore) load() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &s.data)
}

func (s *fileStore) persistLocked() error {
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, bytes, 0o644); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}
