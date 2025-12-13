package documentstore

import "errors"

var (
	ErrDocumentNotFound         = errors.New("document not found")
	ErrCollectionAlreadyExists  = errors.New("collection already exists")
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrUnsupportedDocumentField = errors.New("unsupported document field")
	ErrInvalidPrimaryKey        = errors.New("invalid primary key")
)
