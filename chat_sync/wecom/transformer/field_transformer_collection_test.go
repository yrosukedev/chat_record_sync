package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"io"
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

func TestFieldTransformerCollection_Transform_MultipleTransformers(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	transformer1 := NewMockFieldTransformer(ctrl)
	transformer2 := NewMockFieldTransformer(ctrl)
	transformer3 := NewMockFieldTransformer(ctrl)
	collection := NewFieldTransformerCollection([]FieldTransformer{transformer1, transformer2, transformer3})
	wecomRecord := &wecom.ChatRecord{}
	chatRecordStep1 := &business.ChatRecord{}
	chatRecordStep2 := &business.ChatRecord{}
	expectedChatRecord := &business.ChatRecord{}

	transformer1.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Nil()).Return(chatRecordStep1, nil).Times(1)
	transformer2.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(chatRecordStep1)).Return(chatRecordStep2, nil).Times(1)
	transformer3.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(chatRecordStep2)).Return(expectedChatRecord, nil).Times(1)

	// When
	chatRecord, err := collection.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestFieldTransformerCollection_Transform_Error(t *testing.T) {
	// 3 transformers, 2nd one returns error, the error should be ignored

	// Given
	ctrl := gomock.NewController(t)
	transformer1 := NewMockFieldTransformer(ctrl)
	transformer2 := NewMockFieldTransformer(ctrl)
	transformer3 := NewMockFieldTransformer(ctrl)
	collection := NewFieldTransformerCollection([]FieldTransformer{transformer1, transformer2, transformer3})
	wecomRecord := &wecom.ChatRecord{}
	chatRecordStep1 := &business.ChatRecord{}
	expectedChatRecord := &business.ChatRecord{}

	transformer1.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Nil()).Return(chatRecordStep1, nil).Times(1)
	transformer2.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(chatRecordStep1)).Return(nil, io.EOF).Times(1)
	transformer3.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(chatRecordStep1)).Return(expectedChatRecord, nil).Times(1)

	// When
	chatRecord, err := collection.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
