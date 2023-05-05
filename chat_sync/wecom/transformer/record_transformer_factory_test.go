package transformer

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestRecordTransformerFactory_Transform_TransformerNotFoundForMessageType(t *testing.T) {
	// if the transformer for the message type is not found, the default transformer should be used.
	// both the default transformer and the transformer for the message type should be mocked.

	// Given
	ctrl := gomock.NewController(t)
	defaultTransformer := NewMockFieldTransformer(ctrl)
	fieldATransformer := NewMockFieldTransformer(ctrl)
	fieldBTransformer := NewMockFieldTransformer(ctrl)
	messageType := "message type x"
	messageTypeToTransformer := map[string]FieldTransformer{
		fmt.Sprintf("::message type not %v::", messageType):         fieldATransformer,
		fmt.Sprintf("::message type not %v neither::", messageType): fieldBTransformer,
	}
	factory := NewRecordTransformerFactory(messageTypeToTransformer, defaultTransformer)
	wecomRecord := &wecom.ChatRecord{
		MsgType: messageType,
	}
	expectedChatRecord := &business.ChatRecord{
		MsgType: messageType,
	}

	defaultTransformer.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Nil()).Return(expectedChatRecord, nil).Times(1)
	fieldATransformer.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(expectedChatRecord)).Return(expectedChatRecord, nil).Times(0)
	fieldBTransformer.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Eq(expectedChatRecord)).Return(expectedChatRecord, nil).Times(0)

	// When
	chatRecord, err := factory.Transform(wecomRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestRecordTransformerFactory_Transform_TransformerFoundForMessageType(t *testing.T) {
	// if the transformer for the message type is found, the transformer for the message type should be used.
	// both the default transformer and the transformer for the message type should be mocked.

	// Given
	ctrl := gomock.NewController(t)
	defaultTransformer := NewMockFieldTransformer(ctrl)
	fieldATransformer := NewMockFieldTransformer(ctrl)
	fieldBTransformer := NewMockFieldTransformer(ctrl)
	fieldXTransformer := NewMockFieldTransformer(ctrl)
	messageType := "message type x"
	messageTypeToTransformer := map[string]FieldTransformer{
		fmt.Sprintf("::message type not %v::", messageType):         fieldATransformer,
		fmt.Sprintf("::message type not %v neither::", messageType): fieldBTransformer,
		messageType: fieldXTransformer,
	}
	factory := NewRecordTransformerFactory(messageTypeToTransformer, defaultTransformer)
	wecomRecord := &wecom.ChatRecord{
		MsgType: messageType,
	}
	expectedChatRecord := &business.ChatRecord{
		MsgType: messageType,
	}

	defaultTransformer.EXPECT().Transform(gomock.Any(), gomock.Any()).Return(expectedChatRecord, nil).Times(0)
	fieldATransformer.EXPECT().Transform(gomock.Any(), gomock.Any()).Return(expectedChatRecord, nil).Times(0)
	fieldBTransformer.EXPECT().Transform(gomock.Any(), gomock.Any()).Return(expectedChatRecord, nil).Times(0)
	fieldXTransformer.EXPECT().Transform(gomock.Eq(wecomRecord), gomock.Nil()).Return(expectedChatRecord, nil).Times(1)

	// When
	chatRecord, err := factory.Transform(wecomRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
