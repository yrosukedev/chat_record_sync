package wecom

import (
	"fmt"
)

type TransformerEmptyContentError struct {
	record *ChatRecord
}

func NewTransformerEmptyContentError(record *ChatRecord) *TransformerEmptyContentError {
	return &TransformerEmptyContentError{
		record: record,
	}
}

func (t *TransformerEmptyContentError) Error() string {
	return fmt.Sprintf("transformer error: content should not be empty, record: %#v", t.record)
}
