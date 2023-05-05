package transformer

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestFieldTransformerCollection_Transform_ZeroTransformer(t *testing.T) {
	// Given
	collection := NewFieldTransformerCollection(nil)
	wecomRecord := &wecom.ChatRecord{}
	var expectedChatRecord *business.ChatRecord

	// When
	chatRecord, err := collection.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
