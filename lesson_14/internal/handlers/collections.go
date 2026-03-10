package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
