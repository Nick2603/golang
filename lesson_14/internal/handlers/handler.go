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

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
