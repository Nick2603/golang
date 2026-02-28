package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Storage interface {
	PutDocument(ctx context.Context, collectionName string, document map[string]any) error
	GetDocument(ctx context.Context, collectionName string, filter map[string]any) (map[string]any, error)
	ListDocuments(ctx context.Context, collectionName string, filter map[string]any) ([]map[string]any, error)
	DeleteDocument(ctx context.Context, collectionName string, filter map[string]any) (int64, error)
	CreateCollection(ctx context.Context, collectionName string) error
	ListCollections(ctx context.Context) ([]string, error)
	DeleteCollection(ctx context.Context, collectionName string) error
	CreateIndex(ctx context.Context, collectionName string, fieldName string, unique bool, sparse bool) (string, error)
	DeleteIndex(ctx context.Context, collectionName string, indexName string) error
}

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
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

	if err := h.storage.PutDocument(r.Context(), req.CollectionName, req.Document); err != nil {
		http.Error(w, fmt.Sprintf("failed to put document: %v", err), http.StatusInternalServerError)
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

	doc, err := h.storage.GetDocument(r.Context(), req.CollectionName, req.Filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find document: %v", err), http.StatusInternalServerError)
		return
	}

	if doc == nil {
		respondJSON(w, GetDocumentResponse{Ok: false})
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

	documents, err := h.storage.ListDocuments(r.Context(), req.CollectionName, req.Filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find documents: %v", err), http.StatusInternalServerError)
		return
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

	deletedCount, err := h.storage.DeleteDocument(r.Context(), req.CollectionName, req.Filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete document: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, DeleteDocumentResponse{Ok: deletedCount > 0, DeletedCount: deletedCount})
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

	if err := h.storage.CreateCollection(r.Context(), req.CollectionName); err != nil {
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
	collections, err := h.storage.ListCollections(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list collections: %v", err), http.StatusInternalServerError)
		return
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

	if err := h.storage.DeleteCollection(r.Context(), req.CollectionName); err != nil {
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

	indexName, err := h.storage.CreateIndex(r.Context(), req.CollectionName, req.FieldName, req.Unique, req.Sparse)
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

	if err := h.storage.DeleteIndex(r.Context(), req.CollectionName, req.IndexName); err != nil {
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
