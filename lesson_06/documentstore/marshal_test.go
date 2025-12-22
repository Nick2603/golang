package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Age    int64  `json:"age"`
	Active bool   `json:"active"`
}

func TestMarshalDocument(t *testing.T) {
	t.Run("marshals struct successfully", func(t *testing.T) {
		input := TestStruct{
			ID:     "user:1",
			Name:   "Alice",
			Age:    25,
			Active: true,
		}

		doc, err := MarshalDocument(&input)
		assert.NoError(t, err)
		require.NotNil(t, doc)

		assert.Equal(t, "user:1", doc.Fields["id"].Value)
		assert.Equal(t, "Alice", doc.Fields["name"].Value)
		assert.Equal(t, int64(25), doc.Fields["age"].Value)
		assert.Equal(t, true, doc.Fields["active"].Value)
	})

	t.Run("marshals struct without pointer", func(t *testing.T) {
		input := TestStruct{ID: "user:2", Name: "Bob", Age: 30, Active: false}

		doc, err := MarshalDocument(input)
		assert.NoError(t, err)
		require.NotNil(t, doc)
	})

	t.Run("returns error for non-struct input", func(t *testing.T) {
		input := "not a struct"

		doc, err := MarshalDocument(input)
		assert.ErrorIs(t, err, ErrUnsupportedDocumentField)
		assert.Nil(t, doc)
	})
}

func TestUnmarshalDocument(t *testing.T) {
	t.Run("unmarshals document successfully", func(t *testing.T) {
		doc := &Document{
			Fields: map[string]DocumentField{
				"id":     {Type: DocumentFieldTypeString, Value: "user:1"},
				"name":   {Type: DocumentFieldTypeString, Value: "Alice"},
				"age":    {Type: DocumentFieldTypeNumber, Value: int64(25)},
				"active": {Type: DocumentFieldTypeBool, Value: true},
			},
		}

		var output TestStruct
		err := UnmarshalDocument(doc, &output)
		assert.NoError(t, err)

		assert.Equal(t, "user:1", output.ID)
		assert.Equal(t, "Alice", output.Name)
		assert.Equal(t, int64(25), output.Age)
		assert.Equal(t, true, output.Active)
	})

	t.Run("returns error when output is not a pointer", func(t *testing.T) {
		doc := &Document{Fields: map[string]DocumentField{}}

		var output TestStruct
		err := UnmarshalDocument(doc, output)
		assert.ErrorIs(t, err, ErrUnsupportedDocumentField)
	})

	t.Run("returns error when output is not a struct pointer", func(t *testing.T) {
		doc := &Document{Fields: map[string]DocumentField{}}

		var output string
		err := UnmarshalDocument(doc, &output)
		assert.ErrorIs(t, err, ErrUnsupportedDocumentField)
	})
}
