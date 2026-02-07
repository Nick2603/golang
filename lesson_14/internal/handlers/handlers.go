package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	client *mongo.Client
	dbName string
}

func NewHandler(client *mongo.Client, dbName string) *Handler {
	return &Handler{
		client: client,
		dbName: dbName,
	}
}

type PutDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Document       map[string]any `json:"document"`
}

type PutDocumentResponse struct {
	Ok bool `json:"ok"`
}

func (h *Handler) HandlePutDocument(w http.ResponseWriter, r *http.Request) {
	var req PutDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	if len(req.Document) == 0 {
		http.Error(w, "document is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	var filter bson.M
	if id, ok := req.Document["_id"]; ok {
		filter = bson.M{"_id": id}
	} else if id, ok := req.Document["id"]; ok {
		filter = bson.M{"id": id}
	} else {
		_, err := coll.InsertOne(r.Context(), req.Document)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to insert document: %v", err), http.StatusInternalServerError)
			return
		}
		respondJSON(w, PutDocumentResponse{Ok: true})
		return
	}

	update := bson.M{"$set": req.Document}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update document: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, PutDocumentResponse{Ok: true})
}

type GetDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter"`
}

type GetDocumentResponse struct {
	Ok       bool           `json:"ok"`
	Document map[string]any `json:"document,omitempty"`
}

func (h *Handler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	var req GetDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	if len(req.Filter) == 0 {
		http.Error(w, "filter is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	var doc map[string]any
	err := coll.FindOne(r.Context(), req.Filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondJSON(w, GetDocumentResponse{Ok: false})
			return
		}
		http.Error(w, fmt.Sprintf("failed to find document: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, GetDocumentResponse{Ok: true, Document: doc})
}

type ListDocumentsRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter,omitempty"`
}

type ListDocumentsResponse struct {
	Ok        bool             `json:"ok"`
	Documents []map[string]any `json:"documents"`
}

func (h *Handler) HandleListDocuments(w http.ResponseWriter, r *http.Request) {
	var req ListDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	filter := bson.M{}
	if len(req.Filter) > 0 {
		filter = req.Filter
	}

	cursor, err := coll.Find(r.Context(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find documents: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var documents []map[string]any
	for cursor.Next(r.Context()) {
		var doc map[string]any
		if err := cursor.Decode(&doc); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode document: %v", err), http.StatusInternalServerError)
			return
		}
		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, fmt.Sprintf("cursor error: %v", err), http.StatusInternalServerError)
		return
	}

	if documents == nil {
		documents = []map[string]any{}
	}

	respondJSON(w, ListDocumentsResponse{Ok: true, Documents: documents})
}

type DeleteDocumentRequest struct {
	CollectionName string         `json:"collection_name"`
	Filter         map[string]any `json:"filter"`
}

type DeleteDocumentResponse struct {
	Ok           bool  `json:"ok"`
	DeletedCount int64 `json:"deleted_count"`
}

func (h *Handler) HandleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	var req DeleteDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	if len(req.Filter) == 0 {
		http.Error(w, "filter is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	result, err := coll.DeleteOne(r.Context(), req.Filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete document: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, DeleteDocumentResponse{
		Ok:           result.DeletedCount > 0,
		DeletedCount: result.DeletedCount,
	})
}

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type CreateCollectionResponse struct {
	Ok bool `json:"ok"`
}

func (h *Handler) HandleCreateCollection(w http.ResponseWriter, r *http.Request) {
	var req CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	db := h.client.Database(h.dbName)

	err := db.CreateCollection(r.Context(), req.CollectionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create collection: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, CreateCollectionResponse{Ok: true})
}

type ListCollectionsRequest struct{}

type ListCollectionsResponse struct {
	Ok          bool     `json:"ok"`
	Collections []string `json:"collections"`
}

func (h *Handler) HandleListCollections(w http.ResponseWriter, r *http.Request) {
	db := h.client.Database(h.dbName)

	collections, err := db.ListCollectionNames(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list collections: %v", err), http.StatusInternalServerError)
		return
	}

	if collections == nil {
		collections = []string{}
	}

	respondJSON(w, ListCollectionsResponse{Ok: true, Collections: collections})
}

type DeleteCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type DeleteCollectionResponse struct {
	Ok bool `json:"ok"`
}

func (h *Handler) HandleDeleteCollection(w http.ResponseWriter, r *http.Request) {
	var req DeleteCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	err := coll.Drop(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete collection: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, DeleteCollectionResponse{Ok: true})
}

type CreateIndexRequest struct {
	CollectionName string `json:"collection_name"`
	FieldName      string `json:"field_name"`
	Unique         bool   `json:"unique,omitempty"`
	Sparse         bool   `json:"sparse,omitempty"`
}

type CreateIndexResponse struct {
	Ok        bool   `json:"ok"`
	IndexName string `json:"index_name,omitempty"`
}

func (h *Handler) HandleCreateIndex(w http.ResponseWriter, r *http.Request) {
	var req CreateIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	if req.FieldName == "" {
		http.Error(w, "field_name is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: req.FieldName, Value: 1}}, // 1 for ascending
	}

	if req.Unique || req.Sparse {
		opts := options.Index()
		if req.Unique {
			opts.SetUnique(true)
		}
		if req.Sparse {
			opts.SetSparse(true)
		}
		indexModel.Options = opts
	}

	indexName, err := coll.Indexes().CreateOne(r.Context(), indexModel)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create index: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, CreateIndexResponse{Ok: true, IndexName: indexName})
}

type DeleteIndexRequest struct {
	CollectionName string `json:"collection_name"`
	IndexName      string `json:"index_name"`
}

type DeleteIndexResponse struct {
	Ok bool `json:"ok"`
}

func (h *Handler) HandleDeleteIndex(w http.ResponseWriter, r *http.Request) {
	var req DeleteIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	if req.IndexName == "" {
		http.Error(w, "index_name is required", http.StatusBadRequest)
		return
	}

	coll := h.client.Database(h.dbName).Collection(req.CollectionName)

	_, err := coll.Indexes().DropOne(r.Context(), req.IndexName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete index: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, DeleteIndexResponse{Ok: true})
}

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
