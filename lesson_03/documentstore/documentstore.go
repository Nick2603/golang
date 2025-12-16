package documentstore

import "fmt"

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]*Document{}

func Put(doc *Document) error {
	if doc == nil || doc.Fields == nil {
		return fmt.Errorf("document or fields cannot be nil")
	}

	keyField, exists := doc.Fields["key"]

	if !exists {
		return fmt.Errorf("document must contain 'key' field")
	}

	if keyField.Type != DocumentFieldTypeString {
		return fmt.Errorf("'key' field must be of type string")
	}

	key, ok := keyField.Value.(string)

	if !ok || key == "" {
		return fmt.Errorf("'key' field value must be a non-empty string")
	}

	documents[key] = doc

	return nil
}

func Get(key string) (*Document, bool) {
	doc, exists := documents[key]
	if exists {
		return doc, true
	}

	return nil, false
}

func Delete(key string) bool {
	if _, exists := documents[key]; exists {
		delete(documents, key)

		return true
	}

	return false
}

func List() []*Document {
	result := make([]*Document, 0, len(documents))

	for _, doc := range documents {
		result = append(result, doc)
	}

	return result
}
