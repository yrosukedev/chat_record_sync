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
		MsgType: "text",
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

func TestTextContentFieldTransformer_Transform_dontChangeInputs(t *testing.T) {
	// Given
	transformer := NewTextContentFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		MsgType: "text",
		Text: &wecom.TextMessage{
			Content: "Hello, there!",
		},
	}
	chatRecord := &business.ChatRecord{
		Content: "::whatever content::",
	}
	expectedChatRecord := &business.ChatRecord{
		Content: "Hello, there!",
	}

	// When
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)
		assert.Equal(t, "::whatever content::", chatRecord.Content)
	}
}

func TestTextContentFieldTransformer_Transform_mismatchedMessageType(t *testing.T) {
	// if the message type is not "text", the transformer should panic.

	// Given
	transformer := NewTextContentFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		MsgType: "image",
		Text: &wecom.TextMessage{
			Content: "Hello, there!",
		},
	}

	// Then
	assert.Panics(t, func() {
		// When
		_, _ = transformer.Transform(wecomRecord, nil)
	})
}

func TestTextContentFieldTransformer_Transform_nilTextMessage(t *testing.T) {
	// if the message type is "text", but the Text field is nil, treat it as an empty string.

	// Given
	transformer := NewTextContentFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		MsgType: "text",
	}
	expectedChatRecord := &business.ChatRecord{
		Content: "",
	}

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

