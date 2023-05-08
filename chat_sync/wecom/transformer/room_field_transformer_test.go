package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestRoomFieldTransformer_Transform_NilChatRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	nameFetcher := NewMockNameFetcher(ctrl)
	transformer := NewRoomFieldTransformer(nameFetcher)
	wecomRecord := &wecom.ChatRecord{
		RoomID: "123",
	}
	expectedChatRecord := &business.ChatRecord{
		Room: &business.Room{
			RoomId: "123",
			Name:   "room123",
		},
	}

	nameFetcher.EXPECT().FetchName(gomock.Eq("123")).Return("room123", nil).Times(1)

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestRoomFieldTransformer_Transform_WecomRecordCantBeNil(t *testing.T) {
	// if the wecomRecord is nil, return error

	// Given
	ctrl := gomock.NewController(t)
	nameFetcher := NewMockNameFetcher(ctrl)
	transformer := NewRoomFieldTransformer(nameFetcher)

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}


