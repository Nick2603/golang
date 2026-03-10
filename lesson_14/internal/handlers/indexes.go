package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
