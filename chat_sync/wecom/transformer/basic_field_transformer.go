package transformer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"time"
)

type BasicFieldTransformer struct {
}

func NewBasicFieldTransformer() *BasicFieldTransformer {
	return &BasicFieldTransformer{}
}

func (t *BasicFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if chatRecord == nil {
		chatRecord = &business.ChatRecord{}
	}

	chatRecord.MsgId = wecomRecord.MsgID
	chatRecord.MsgTime = time.UnixMilli(wecomRecord.MsgTime)
	chatRecord.Action = wecomRecord.Action
	chatRecord.MsgType = wecomRecord.MsgType
	chatRecord.RoomId = wecomRecord.RoomID

	return chatRecord, nil
}
