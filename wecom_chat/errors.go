package wecom_chat

import "fmt"

func NewTransformerErrorMessageTypeMissMatched(expectedType string, actualType string) error {
	return fmt.Errorf("transformer error: message type missmatched, expected: %v, actual: %v", expectedType, actualType)
}
