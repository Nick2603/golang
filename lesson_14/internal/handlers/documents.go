package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
