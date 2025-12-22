package documentstore

import "log/slog"

type Store struct {
	collections map[string]*Collection
	logger      *slog.Logger
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

	if _, exists := s.collections[name]; exists {
		s.logger.Warn("collection already exists", "collection", name)
		return nil, ErrCollectionAlreadyExists
	}

	coll := NewCollection(*cfg)
	coll.logger = s.logger
	s.collections[name] = coll

	s.logger.Info("collection created", "collection", name, "primary_key", cfg.PrimaryKey)
	return coll, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	coll, ok := s.collections[name]
	if !ok {
		s.logger.Warn("collection not found", "collection", name)
		return nil, ErrCollectionNotFound
	}
	return coll, nil
}

func (s *Store) DeleteCollection(name string) error {
	if _, ok := s.collections[name]; !ok {
		s.logger.Warn("failed to delete collection: not found", "collection", name)
		return ErrCollectionNotFound
	}

	delete(s.collections, name)
	s.logger.Info("collection deleted", "collection", name)
	return nil
}
