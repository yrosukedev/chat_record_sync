package transformer

import (
	"github.com/golang/mock/gomock"
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

func TestFieldTransformerCollection_Transform_OneTransformer(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	transformer := NewMockFieldTransformer(ctrl)
	collection := NewFieldTransformerCollection([]FieldTransformer{transformer})
	wecomRecord := &wecom.ChatRecord{}
	expectedChatRecord := &business.ChatRecord{}

	transformer.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Nil()).Return(expectedChatRecord, nil).Times(1)

	// When
	chatRecord, err := collection.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
