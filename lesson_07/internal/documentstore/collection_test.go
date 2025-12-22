package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Put(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		wantErr error
	}{
		{
			name: "successfully adds document",
			doc: Document{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
					"name": {Type: DocumentFieldTypeString, Value: "Alice"},
				},
			},
			wantErr: nil,
		},
		{
			name: "returns error when fields are nil",
			doc: Document{
				Fields: nil,
			},
			wantErr: ErrNilValue,
		},
		{
			name: "returns error when primary key is missing",
			doc: Document{
				Fields: map[string]DocumentField{
					"name": {Type: DocumentFieldTypeString, Value: "Test"},
				},
			},
			wantErr: ErrInvalidPrimaryKey,
		},
		{
			name: "returns error when primary key is not string type",
			doc: Document{
				Fields: map[string]DocumentField{
					"id": {Type: DocumentFieldTypeNumber, Value: 123},
				},
			},
			wantErr: ErrInvalidPrimaryKey,
		},
		{
			name: "returns error when primary key value is empty",
			doc: Document{
				Fields: map[string]DocumentField{
					"id": {Type: DocumentFieldTypeString, Value: ""},
				},
			},
			wantErr: ErrInvalidPrimaryKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

			err := coll.Put(tt.doc)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollection_Get(t *testing.T) {
	tests := []struct {
		name      string
		setupDocs []Document
		getKey    string
		wantDoc   *Document
		wantErr   error
	}{
		{
			name: "returns document when exists",
			setupDocs: []Document{
				{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
						"name": {Type: DocumentFieldTypeString, Value: "Alice"},
					},
				},
			},
			getKey: "user:1",
			wantDoc: &Document{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
					"name": {Type: DocumentFieldTypeString, Value: "Alice"},
				},
			},
			wantErr: nil,
		},
		{
			name:      "returns error when document not found",
			setupDocs: []Document{},
			getKey:    "user:999",
			wantDoc:   nil,
			wantErr:   ErrDocumentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

			// Setup
			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

			// Test
			retrieved, err := coll.Get(tt.getKey)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, retrieved)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, retrieved)
				assert.Equal(t, tt.wantDoc.Fields["name"].Value, retrieved.Fields["name"].Value)
			}
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	tests := []struct {
		name      string
		setupDocs []Document
		deleteKey string
		wantErr   error
	}{
		{
			name: "deletes existing document",
			setupDocs: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "user:1"},
					},
				},
			},
			deleteKey: "user:1",
			wantErr:   nil,
		},
		{
			name:      "returns error when document not found",
			setupDocs: []Document{},
			deleteKey: "user:999",
			wantErr:   ErrDocumentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

			// Setup
			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

			// Test
			err := coll.Delete(tt.deleteKey)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Verify deletion
				_, getErr := coll.Get(tt.deleteKey)
				assert.ErrorIs(t, getErr, ErrDocumentNotFound)
			}
		})
	}
}

func TestCollection_List(t *testing.T) {
	tests := []struct {
		name      string
		setupDocs []Document
		wantCount int
	}{
		{
			name: "returns all documents",
			setupDocs: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "user:1"},
					},
				},
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "user:2"},
					},
				},
			},
			wantCount: 2,
		},
		{
			name:      "returns empty slice when collection is empty",
			setupDocs: []Document{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

			// Setup
			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

			// Test
			docs := coll.List()

			assert.Len(t, docs, tt.wantCount)
			assert.NotNil(t, docs)
		})
	}
}
