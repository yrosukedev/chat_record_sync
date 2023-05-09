package formatter

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/consts"
	"testing"
	"time"
)

func TestBitableFieldsFormatter_Format(t *testing.T) {
	// Given
	formatter := NewBitableFieldsFormatter()

	chatRecord := &business.ChatRecord{
		MsgId:   "::whatever MsgId::",
		Action:  "::whatever Action::",
		MsgType: "::whatever MsgType::",
		MsgTime: time.UnixMilli(1610000000000),
		Content: "::whatever Content::",
		From: &business.User{
			UserId: "::whatever UserId::",
			Name:   "::whatever UserName::",
		},
		To: []*business.User{
			{
				UserId: "::whatever Contact Id A::",
				Name:   "::whatever Contact Name A::",
			},
			{
				UserId: "::whatever Contact Id B::",
				Name:   "::whatever Contact Name B::",
			},
		},
		Room: &business.Room{
			RoomId: "::whatever RoomId::",
			Name:   "::whatever RoomName::",
		},
	}
	expectedFields := map[string]interface{}{
		consts.BitableFieldChatRecordMsgId:   "::whatever MsgId::",
		consts.BitableFieldChatRecordAction:  "::whatever Action::",
		consts.BitableFieldChatRecordMsgType: "::whatever MsgType::",
		consts.BitableFieldChatRecordMsgTime: int64(1610000000000),
		consts.BitableFieldChatRecordContent: "::whatever Content::",
		consts.BitableFieldChatRecordFrom:    "::whatever UserName::(ID:::whatever UserId::)",
		consts.BitableFieldChatRecordTo:      "::whatever Contact Name A::(ID:::whatever Contact Id A::),::whatever Contact Name B::(ID:::whatever Contact Id B::)",
		consts.BitableFieldChatRecordRoomId:  "::whatever RoomId::",
	}

	// When
	fields, err := formatter.Format(chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedFields, fields)
	}
}
