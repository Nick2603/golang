package documentstore

import "fmt"

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		cfg:       cfg,
		documents: make(map[string]Document),
	}
}

func (c *Collection) Put(doc Document) error {
	if doc.Fields == nil {
		return fmt.Errorf("document fields cannot be nil")
	}

	pkField, exists := doc.Fields[c.cfg.PrimaryKey]
	if !exists {
		return fmt.Errorf("document must contain primary key field '%s'", c.cfg.PrimaryKey)
	}

	if pkField.Type != DocumentFieldTypeString {
		return fmt.Errorf("primary key field '%s' must be of type string", c.cfg.PrimaryKey)
	}

	key, ok := pkField.Value.(string)
	if !ok || key == "" {
		return fmt.Errorf("primary key value must be a non-empty string")
	}

	c.documents[key] = doc
	return nil
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, exists := c.documents[key]
	if !exists {
		return nil, false
	}

	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	if _, exists := c.documents[key]; exists {
		delete(c.documents, key)
		return true
	}

	return false
}

func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.documents))

	for _, doc := range c.documents {
		result = append(result, doc)
	}

	return result
}
