package documentstore

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, exists := s.collections[name]; exists {
		return nil, ErrCollectionAlreadyExists
	}

	coll := NewCollection(*cfg)
	s.collections[name] = coll
	return coll, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	coll, ok := s.collections[name]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return coll, nil
}

func (s *Store) DeleteCollection(name string) error {
	if _, ok := s.collections[name]; !ok {
		return ErrCollectionNotFound
	}
	delete(s.collections, name)
	return nil
}
