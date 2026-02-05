package documentstore

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value any
}

type Document struct {
	Fields map[string]DocumentField
}
