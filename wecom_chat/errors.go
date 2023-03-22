package wecom_chat

import (
	"fmt"
)

func NewTransformerErrorMessageTypeMissMatched(expectedType string, actualType string) error {
	return fmt.Errorf("transformer error: message type missmatched, expected: %v, actual: %v", expectedType, actualType)
}

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
