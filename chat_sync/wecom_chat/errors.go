package wecom_chat

import (
	"fmt"
)

type TransformerEmptyContentError struct {
	record *WeComChatRecord
}

func NewTransformerEmptyContentError(record *WeComChatRecord) *TransformerEmptyContentError {
	return &TransformerEmptyContentError{
		record: record,
	}
}

func (t *TransformerEmptyContentError) Error() string {
	return fmt.Sprintf("transformer error: content should not be empty, record: %#v", t.record)
}
