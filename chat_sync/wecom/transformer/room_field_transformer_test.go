package transformer

import (
	"errors"
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

func TestRoomFieldTransformer_Transform_DontChangeInputs(t *testing.T) {
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
	chatRecord := &business.ChatRecord{
		Room: &business.Room{
			RoomId: "::whatever room id that can't be changed::",
			Name:   "::whatever room name that can't be changed::",
		},
	}

	nameFetcher.EXPECT().FetchName(gomock.Eq("123")).Return("room123", nil).Times(1)

	// When
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)
		assert.Equal(t, "::whatever room id that can't be changed::", chatRecord.Room.RoomId, "chatRecord.Room.RoomId should not be changed")
		assert.Equal(t, "::whatever room name that can't be changed::", chatRecord.Room.Name, "chatRecord.Room.Name should not be changed")
	}
}

func TestRoomFieldTransformer_Transform_FetchNameError(t *testing.T) {
	// if the nameFetcher returns error, only populate the RoomId field and leave the Name field untouched,
	// and return nil error

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
		},
	}

	nameFetcher.EXPECT().FetchName(gomock.Eq("123")).Return("", errors.New("some error")).Times(1)

	// When
	updatedChatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)
	}
}

