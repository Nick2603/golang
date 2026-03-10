package documentstore

import (
	"fmt"
	"log/slog"
	"sync"
)

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	indexes   map[string]*Index
	logger    *slog.Logger
	mu        sync.RWMutex
}

type CollectionConfig struct {
	PrimaryKey string
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		cfg:       cfg,
		documents: make(map[string]Document),
		indexes:   make(map[string]*Index),
		logger:    slog.Default(),
	}
}

func (c *Collection) Put(doc Document) error {
	if doc.Fields == nil {
		c.logger.Error("failed to put document: nil fields")
		return ErrNilValue
	}

	field, ok := doc.Fields[c.cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		c.logger.Error("failed to put document: invalid primary key", "primary_key", c.cfg.PrimaryKey)
		return ErrInvalidPrimaryKey
	}

	key, ok := field.Value.(string)
	if !ok || key == "" {
		c.logger.Error("failed to put document: empty primary key value")
		return ErrInvalidPrimaryKey
	}

	c.mu.Lock()
	oldDoc, exists := c.documents[key]

	if exists {
		c.removeFromIndexes(key, oldDoc)
	}

	c.documents[key] = doc

	c.updateIndexes(key, doc)
	c.mu.Unlock()

	if exists {
		c.logger.Info("document updated", "key", key)
	} else {
		c.logger.Info("document created", "key", key)
	}

	return nil
}

func (c *Collection) Get(key string) (*Document, error) {
	c.mu.RLock()
	doc, ok := c.documents[key]
	c.mu.RUnlock()

	if !ok {
		c.logger.Warn("document not found", "key", key)
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func (c *Collection) Delete(key string) error {
	c.mu.Lock()
	doc, ok := c.documents[key]
	if !ok {
		c.mu.Unlock()
		c.logger.Warn("failed to delete document: not found", "key", key)
		return ErrDocumentNotFound
	}

	c.removeFromIndexes(key, doc)

	delete(c.documents, key)
	c.mu.Unlock()

	c.logger.Info("document deleted", "key", key)
	return nil
}

func (c *Collection) List() []Document {
	c.mu.RLock()
	result := make([]Document, 0, len(c.documents))
	for _, d := range c.documents {
		result = append(result, d)
	}
	c.mu.RUnlock()

	return result
}

func (c *Collection) CreateIndex(fieldName string) error {
	c.mu.Lock()
	if c.indexes == nil {
		c.indexes = make(map[string]*Index)
	}

	if _, exists := c.indexes[fieldName]; exists {
		c.mu.Unlock()
		c.logger.Warn("index already exists", "field", fieldName)
		return fmt.Errorf("index already exists for field: %s", fieldName)
	}

	idx := newIndex(fieldName)

	for primaryKey, doc := range c.documents {
		field, exists := doc.Fields[fieldName]
		if !exists {
			continue
		}

		if field.Type != DocumentFieldTypeString {
			continue
		}

		value, ok := field.Value.(string)
		if !ok {
			continue
		}

		idx.insert(value, primaryKey)
	}

	c.indexes[fieldName] = idx
	c.mu.Unlock()

	c.logger.Info("index created", "field", fieldName, "entries", idx.len())

	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	c.mu.Lock()
	if c.indexes == nil {
		c.indexes = make(map[string]*Index)
	}

	if _, exists := c.indexes[fieldName]; !exists {
		c.mu.Unlock()
		c.logger.Warn("index not found", "field", fieldName)
		return fmt.Errorf("index not found for field: %s", fieldName)
	}

	delete(c.indexes, fieldName)
	c.mu.Unlock()

	c.logger.Info("index deleted", "field", fieldName)

	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.mu.RLock()

	if c.indexes == nil {
		c.mu.RUnlock()
		c.logger.Warn("index not found for query", "field", fieldName)
		return nil, fmt.Errorf("index not found for field: %s", fieldName)
	}

	idx, exists := c.indexes[fieldName]
	if !exists {
		c.mu.RUnlock()
		c.logger.Warn("index not found for query", "field", fieldName)
		return nil, fmt.Errorf("index not found for field: %s", fieldName)
	}

	entries := idx.query(params)

	results := make([]Document, 0, len(entries))
	for _, entry := range entries {
		if doc, ok := c.documents[entry.PrimaryKey]; ok {
			results = append(results, doc)
		}
	}
	c.mu.RUnlock()

	c.logger.Info("query executed", "field", fieldName, "results", len(results))

	return results, nil
}

func (c *Collection) updateIndexes(primaryKey string, doc Document) {
	if c.indexes == nil {
		return
	}

	for fieldName, idx := range c.indexes {
		field, exists := doc.Fields[fieldName]
		if !exists || field.Type != DocumentFieldTypeString {
			continue
		}

		value, ok := field.Value.(string)
		if !ok {
			continue
		}

		idx.insert(value, primaryKey)
	}
}

func (c *Collection) removeFromIndexes(primaryKey string, doc Document) {
	if c.indexes == nil {
		return
	}

	for fieldName, idx := range c.indexes {
		field, exists := doc.Fields[fieldName]
		if !exists || field.Type != DocumentFieldTypeString {
			continue
		}

		value, ok := field.Value.(string)
		if !ok {
			continue
		}

		idx.delete(value, primaryKey)
	}
}

func (c *Collection) GetConfig() CollectionConfig {
	return c.cfg
}
