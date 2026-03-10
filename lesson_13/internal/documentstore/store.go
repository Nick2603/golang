package documentstore

import (
	"log/slog"
	"sync"
)

type Store struct {
	collections map[string]*Collection
	logger      *slog.Logger
	mu          sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
		logger:      slog.Default(),
	}
}

func NewStoreWithLogger(logger *slog.Logger) *Store {
	return &Store{
		collections: make(map[string]*Collection),
		logger:      logger,
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if cfg == nil {
		s.logger.Error("failed to create collection: nil config", "collection", name)
		return nil, ErrNilValue
	}

	s.mu.Lock()
	if _, exists := s.collections[name]; exists {
		s.mu.Unlock()
		s.logger.Warn("collection already exists", "collection", name)
		return nil, ErrCollectionAlreadyExists
	}

	coll := NewCollection(*cfg)
	coll.logger = s.logger
	s.collections[name] = coll
	s.mu.Unlock()

	s.logger.Info("collection created", "collection", name, "primary_key", cfg.PrimaryKey)
	return coll, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	s.mu.RLock()
	coll, ok := s.collections[name]
	s.mu.RUnlock()

	if !ok {
		s.logger.Warn("collection not found", "collection", name)
		return nil, ErrCollectionNotFound
	}
	return coll, nil
}

func (s *Store) DeleteCollection(name string) error {
	s.mu.Lock()
	if _, ok := s.collections[name]; !ok {
		s.mu.Unlock()
		s.logger.Warn("failed to delete collection: not found", "collection", name)
		return ErrCollectionNotFound
	}

	delete(s.collections, name)
	s.mu.Unlock()

	s.logger.Info("collection deleted", "collection", name)
	return nil
}
