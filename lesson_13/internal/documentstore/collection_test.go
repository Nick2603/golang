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

			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

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

			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

			err := coll.Delete(tt.deleteKey)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

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

			for _, doc := range tt.setupDocs {
				coll.Put(doc)
			}

			docs := coll.List()

			assert.Len(t, docs, tt.wantCount)
			assert.NotNil(t, docs)
		})
	}
}

func TestCollection_CreateIndex(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Collection)
		fieldName string
		wantErr   bool
	}{
		{
			name: "creates index successfully",
			setupFunc: func(c *Collection) {
				c.Put(Document{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
						"name": {Type: DocumentFieldTypeString, Value: "Alice"},
					},
				})
			},
			fieldName: "name",
			wantErr:   false,
		},
		{
			name: "returns error when index already exists",
			setupFunc: func(c *Collection) {
				c.CreateIndex("name")
			},
			fieldName: "name",
			wantErr:   true,
		},
		{
			name: "creates index for empty collection",
			setupFunc: func(c *Collection) {
			},
			fieldName: "name",
			wantErr:   false,
		},
		{
			name: "indexes only string fields",
			setupFunc: func(c *Collection) {
				c.Put(Document{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
						"name": {Type: DocumentFieldTypeString, Value: "Alice"},
						"age":  {Type: DocumentFieldTypeNumber, Value: int64(25)},
					},
				})
			},
			fieldName: "name",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})
			tt.setupFunc(coll)

			err := coll.CreateIndex(tt.fieldName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, coll.indexes[tt.fieldName])
			}
		})
	}
}

func TestCollection_DeleteIndex(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Collection)
		fieldName string
		wantErr   bool
	}{
		{
			name: "deletes existing index",
			setupFunc: func(c *Collection) {
				c.CreateIndex("name")
			},
			fieldName: "name",
			wantErr:   false,
		},
		{
			name: "returns error when index does not exist",
			setupFunc: func(c *Collection) {
			},
			fieldName: "name",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(CollectionConfig{PrimaryKey: "id"})
			tt.setupFunc(coll)

			err := coll.DeleteIndex(tt.fieldName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, coll.indexes[tt.fieldName])
			}
		})
	}
}

func TestCollection_Query(t *testing.T) {
	setupCollection := func() *Collection {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		docs := []Document{
			{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
					"name": {Type: DocumentFieldTypeString, Value: "Alice"},
				},
			},
			{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:2"},
					"name": {Type: DocumentFieldTypeString, Value: "Bob"},
				},
			},
			{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:3"},
					"name": {Type: DocumentFieldTypeString, Value: "Charlie"},
				},
			},
			{
				Fields: map[string]DocumentField{
					"id":   {Type: DocumentFieldTypeString, Value: "user:4"},
					"name": {Type: DocumentFieldTypeString, Value: "Diana"},
				},
			},
		}

		for _, doc := range docs {
			coll.Put(doc)
		}

		return coll
	}

	t.Run("returns error when index does not exist", func(t *testing.T) {
		coll := setupCollection()

		_, err := coll.Query("name", QueryParams{})

		assert.Error(t, err)
	})

	t.Run("queries all documents with no filters", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		results, err := coll.Query("name", QueryParams{})

		assert.NoError(t, err)
		assert.Len(t, results, 4)
	})

	t.Run("queries with ascending order", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		results, err := coll.Query("name", QueryParams{Desc: false})

		assert.NoError(t, err)
		require.Len(t, results, 4)
		assert.Equal(t, "Alice", results[0].Fields["name"].Value)
		assert.Equal(t, "Diana", results[3].Fields["name"].Value)
	})

	t.Run("queries with descending order", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		results, err := coll.Query("name", QueryParams{Desc: true})

		assert.NoError(t, err)
		require.Len(t, results, 4)
		assert.Equal(t, "Diana", results[0].Fields["name"].Value)
		assert.Equal(t, "Alice", results[3].Fields["name"].Value)
	})

	t.Run("queries with min value filter", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		minValue := "Charlie"
		results, err := coll.Query("name", QueryParams{
			MinValue: &minValue,
		})

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
		for _, doc := range results {
			assert.GreaterOrEqual(t, doc.Fields["name"].Value.(string), "Charlie")
		}
	})

	t.Run("queries with max value filter", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		maxValue := "Bob"
		results, err := coll.Query("name", QueryParams{
			MaxValue: &maxValue,
		})

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
		for _, doc := range results {
			assert.LessOrEqual(t, doc.Fields["name"].Value.(string), "Bob")
		}
	})

	t.Run("queries with both min and max value filters", func(t *testing.T) {
		coll := setupCollection()
		coll.CreateIndex("name")

		minValue := "Bob"
		maxValue := "Charlie"
		results, err := coll.Query("name", QueryParams{
			MinValue: &minValue,
			MaxValue: &maxValue,
		})

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
		for _, doc := range results {
			name := doc.Fields["name"].Value.(string)
			assert.GreaterOrEqual(t, name, "Bob")
			assert.LessOrEqual(t, name, "Charlie")
		}
	})
}

func TestCollection_IndexMaintenance(t *testing.T) {
	t.Run("index updates when document is added", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})
		coll.CreateIndex("name")

		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		})

		results, err := coll.Query("name", QueryParams{})
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("index updates when document is modified", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})
		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		})

		coll.CreateIndex("name")

		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alicia"},
			},
		})

		results, err := coll.Query("name", QueryParams{})
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Alicia", results[0].Fields["name"].Value)
	})

	t.Run("index updates when document is deleted", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})
		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		})

		coll.CreateIndex("name")
		coll.Delete("user:1")

		results, err := coll.Query("name", QueryParams{})
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})
}

func TestCollection_IndexWithDump(t *testing.T) {
	t.Run("indexes are preserved in dump and restore", func(t *testing.T) {
		originalStore := NewStore()
		coll, _ := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		})
		coll.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:2"},
				"name": {Type: DocumentFieldTypeString, Value: "Bob"},
			},
		})

		coll.CreateIndex("name")

		data, err := originalStore.Dump()
		require.NoError(t, err)

		restoredStore, err := NewStoreFromDump(data)
		require.NoError(t, err)

		restoredColl, err := restoredStore.GetCollection("users")
		require.NoError(t, err)

		results, err := restoredColl.Query("name", QueryParams{})
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})
}
