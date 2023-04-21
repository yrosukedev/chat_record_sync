package transformer

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestTextContentFieldTransformer_Transform_nilChatRecord(t *testing.T) {
	// Given
	transformer := NewTextContentFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		Text: &wecom.TextMessage{
			Content: "Hello, there!",
		},
	}
	expectedChatRecord := &business.ChatRecord{
		Content: "Hello, there!",
	}

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestTextContentFieldTransformer_Transform_wecomRecordCantBeNil(t *testing.T) {
	// Given
	transformer := NewTextContentFieldTransformer()

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}
