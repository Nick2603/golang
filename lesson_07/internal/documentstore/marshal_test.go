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
	tests := []struct {
		name    string
		input   any
		wantErr error
		check   func(*testing.T, *Document)
	}{
		{
			name: "marshals struct successfully",
			input: &TestStruct{
				ID:     "user:1",
				Name:   "Alice",
				Age:    25,
				Active: true,
			},
			wantErr: nil,
			check: func(t *testing.T, doc *Document) {
				assert.Equal(t, "user:1", doc.Fields["id"].Value)
				assert.Equal(t, "Alice", doc.Fields["name"].Value)
				assert.Equal(t, int64(25), doc.Fields["age"].Value)
				assert.Equal(t, true, doc.Fields["active"].Value)
			},
		},
		{
			name: "marshals struct without pointer",
			input: TestStruct{
				ID:     "user:2",
				Name:   "Bob",
				Age:    30,
				Active: false,
			},
			wantErr: nil,
			check: func(t *testing.T, doc *Document) {
				assert.Equal(t, "user:2", doc.Fields["id"].Value)
			},
		},
		{
			name:    "returns error for non-struct input",
			input:   "not a struct",
			wantErr: ErrUnsupportedDocumentField,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := MarshalDocument(tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, doc)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, doc)
				if tt.check != nil {
					tt.check(t, doc)
				}
			}
		})
	}
}

func TestUnmarshalDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     *Document
		output  any
		wantErr error
		check   func(*testing.T, any)
	}{
		{
			name: "unmarshals document successfully",
			doc: &Document{
				Fields: map[string]DocumentField{
					"id":     {Type: DocumentFieldTypeString, Value: "user:1"},
					"name":   {Type: DocumentFieldTypeString, Value: "Alice"},
					"age":    {Type: DocumentFieldTypeNumber, Value: int64(25)},
					"active": {Type: DocumentFieldTypeBool, Value: true},
				},
			},
			output:  &TestStruct{},
			wantErr: nil,
			check: func(t *testing.T, output any) {
				result := output.(*TestStruct)
				assert.Equal(t, "user:1", result.ID)
				assert.Equal(t, "Alice", result.Name)
				assert.Equal(t, int64(25), result.Age)
				assert.Equal(t, true, result.Active)
			},
		},
		{
			name:    "returns error when output is not a pointer",
			doc:     &Document{Fields: map[string]DocumentField{}},
			output:  TestStruct{},
			wantErr: ErrUnsupportedDocumentField,
			check:   nil,
		},
		{
			name:    "returns error when output is not a struct pointer",
			doc:     &Document{Fields: map[string]DocumentField{}},
			output:  new(string),
			wantErr: ErrUnsupportedDocumentField,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalDocument(tt.doc, tt.output)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, tt.output)
				}
			}
		})
	}
}
