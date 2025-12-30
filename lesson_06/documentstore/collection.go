package documentstore

import "log/slog"

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	logger    *slog.Logger
}

type CollectionConfig struct {
	PrimaryKey string
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		cfg:       cfg,
		documents: make(map[string]Document),
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

	_, exists := c.documents[key]
	c.documents[key] = doc

	if exists {
		c.logger.Info("document updated", "key", key)
	} else {
		c.logger.Info("document created", "key", key)
	}

	return nil
}

func (c *Collection) Get(key string) (*Document, error) {
	doc, ok := c.documents[key]
	if !ok {
		c.logger.Warn("document not found", "key", key)
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func (c *Collection) Delete(key string) error {
	if _, ok := c.documents[key]; !ok {
		c.logger.Warn("failed to delete document: not found", "key", key)
		return ErrDocumentNotFound
	}

	delete(c.documents, key)
	c.logger.Info("document deleted", "key", key)
	return nil
}

func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.documents))
	for _, d := range c.documents {
		result = append(result, d)
	}
	return result
}
