package transformer

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
	"time"
)

func TestBasicFieldTransformer_Transform_nilChatRecord(t *testing.T) {
	// Given
	transformer := NewBasicFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		MsgID:   "::whatever msg id::",
		MsgTime: 12345579,
		Action:  "::whatever action::",
		MsgType: "::whatever msg type::",
		RoomID:  "::whatever room id::",
	}
	expectedChatRecord := &business.ChatRecord{
		MsgId:   "::whatever msg id::",
		MsgTime: time.UnixMilli(12345579),
		Action:  "::whatever action::",
		MsgType: "::whatever msg type::",
		Room: &business.Room{
			RoomId: "::whatever room id::",
		},
	}

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestBasicFieldTransformer_Transform_wecomRecordCantBeNil(t *testing.T) {
	// Given
	transformer := NewBasicFieldTransformer()

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}

func TestBasicFieldTransformer_Transform_dontChangeInputs(t *testing.T) {
	// Given
	transformer := NewBasicFieldTransformer()
	wecomRecord := &wecom.ChatRecord{
		MsgID:   "::whatever msg id::",
		MsgTime: 12345579,
		Action:  "::whatever action::",
		MsgType: "::whatever msg type::",
		RoomID:  "::whatever room id::",
	}
	chatRecord := &business.ChatRecord{
		MsgId:   "::whatever msg id can't be changed::",
		MsgTime: time.UnixMilli(9876543),
		Action:  "::whatever action  can't be changed::",
		MsgType: "::whatever msg type  can't be changed::",
		Room: &business.Room{
			RoomId: "::whatever room id can't be changed::",
		},
	}
	expectedChatRecord := &business.ChatRecord{
		MsgId:   "::whatever msg id::",
		MsgTime: time.UnixMilli(12345579),
		Action:  "::whatever action::",
		MsgType: "::whatever msg type::",
		Room: &business.Room{
			RoomId: "::whatever room id::",
		},
	}

	// When
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)

		assert.Equal(t, "::whatever msg id can't be changed::", chatRecord.MsgId)
		assert.Equal(t, time.UnixMilli(9876543), chatRecord.MsgTime)
		assert.Equal(t, "::whatever action  can't be changed::", chatRecord.Action)
		assert.Equal(t, "::whatever msg type  can't be changed::", chatRecord.MsgType)
		assert.Equal(t, "::whatever room id can't be changed::", chatRecord.Room.RoomId)
	}
}
