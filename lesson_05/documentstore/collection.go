package documentstore

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
		return ErrNilValue
	}

	field, ok := doc.Fields[c.cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		return ErrInvalidPrimaryKey
	}

	key, ok := field.Value.(string)
	if !ok || key == "" {
		return ErrInvalidPrimaryKey
	}

	c.documents[key] = doc
	return nil
}

func (c *Collection) Get(key string) (*Document, error) {
	doc, ok := c.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func (c *Collection) Delete(key string) error {
	if _, ok := c.documents[key]; !ok {
		return ErrDocumentNotFound
	}
	delete(c.documents, key)
	return nil
}

func (c *Collection) List() []Document {
	result := make([]Document, 0, len(c.documents))
	for _, d := range c.documents {
		result = append(result, d)
	}
	return result
}
