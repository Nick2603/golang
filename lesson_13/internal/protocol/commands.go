package protocol

import "github.com/Nick2603/golang/lesson_13/internal/documentstore"

const (
	CreateCollectionCommand = "CREATE_COLLECTION"
	DeleteCollectionCommand = "DELETE_COLLECTION"
	GetCollectionCommand    = "GET_COLLECTION"
	ListCollectionsCommand  = "LIST_COLLECTIONS"

	PutDocumentCommand    = "PUT_DOCUMENT"
	GetDocumentCommand    = "GET_DOCUMENT"
	DeleteDocumentCommand = "DELETE_DOCUMENT"
	ListDocumentsCommand  = "LIST_DOCUMENTS"

	CreateIndexCommand = "CREATE_INDEX"
	DeleteIndexCommand = "DELETE_INDEX"
	QueryCommand       = "QUERY"

	PingCommand = "PING"
	QuitCommand = "QUIT"
)

type Request struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
}

type Response struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error"`
}

type CreateCollectionRequest struct {
	Name       string `json:"name"`
	PrimaryKey string `json:"primary_key"`
}

type DeleteCollectionRequest struct {
	Name string `json:"name"`
}

type GetCollectionRequest struct {
	Name string `json:"name"`
}

type CollectionInfo struct {
	Name       string `json:"name"`
	PrimaryKey string `json:"primary_key"`
}

type ListCollectionsResponse struct {
	Collections []CollectionInfo `json:"collections"`
}

type PutDocumentRequest struct {
	Collection string                 `json:"collection"`
	Document   documentstore.Document `json:"document"`
}

type GetDocumentRequest struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}

type GetDocumentResponse struct {
	Document documentstore.Document `json:"document"`
}

type DeleteDocumentRequest struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}

type ListDocumentsRequest struct {
	Collection string `json:"collection"`
}

type ListDocumentsResponse struct {
	Documents []documentstore.Document `json:"documents"`
}

type CreateIndexRequest struct {
	Collection string `json:"collection"`
	FieldName  string `json:"field_name"`
}

type DeleteIndexRequest struct {
	Collection string `json:"collection"`
	FieldName  string `json:"field_name"`
}

type QueryRequest struct {
	Collection string                    `json:"collection"`
	FieldName  string                    `json:"field_name"`
	Params     documentstore.QueryParams `json:"params"`
}

type QueryResponse struct {
	Documents []documentstore.Document `json:"documents"`
}

type PingResponse struct {
	Message string `json:"message"`
}
