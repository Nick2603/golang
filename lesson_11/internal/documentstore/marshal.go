package documentstore

import (
	"fmt"
	"reflect"
)

func MarshalDocument(input any) (*Document, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, ErrUnsupportedDocumentField
	}

	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		name := fieldType.Tag.Get("json")
		if name == "" {
			name = fieldType.Name
		}

		switch fieldValue.Kind() {
		case reflect.String:
			doc.Fields[name] = DocumentField{
				Type:  DocumentFieldTypeString,
				Value: fieldValue.String(),
			}
		case reflect.Int, reflect.Int64:
			doc.Fields[name] = DocumentField{
				Type:  DocumentFieldTypeNumber,
				Value: fieldValue.Int(),
			}
		case reflect.Bool:
			doc.Fields[name] = DocumentField{
				Type:  DocumentFieldTypeBool,
				Value: fieldValue.Bool(),
			}
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedDocumentField, fieldValue.Kind())
		}
	}

	return doc, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	v := reflect.ValueOf(output)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return ErrUnsupportedDocumentField
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		name := fieldType.Tag.Get("json")
		if name == "" {
			name = fieldType.Name
		}

		docField, exists := doc.Fields[name]
		if !exists {
			continue
		}

		fieldValue := v.Field(i)

		switch docField.Type {
		case DocumentFieldTypeString:
			fieldValue.SetString(docField.Value.(string))
		case DocumentFieldTypeNumber:
			fieldValue.SetInt(docField.Value.(int64))
		case DocumentFieldTypeBool:
			fieldValue.SetBool(docField.Value.(bool))
		default:
			return ErrUnsupportedDocumentField
		}
	}

	return nil
}
