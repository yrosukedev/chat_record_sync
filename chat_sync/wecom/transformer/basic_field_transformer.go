package transformer

import (
	"errors"
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
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	updatedChatRecord.MsgId = wecomRecord.MsgID
	updatedChatRecord.MsgTime = time.UnixMilli(wecomRecord.MsgTime)
	updatedChatRecord.Action = wecomRecord.Action
	updatedChatRecord.MsgType = wecomRecord.MsgType
	updatedChatRecord.RoomId = wecomRecord.RoomID

	return updatedChatRecord, nil
}

func (t *BasicFieldTransformer) copyInputIfNeeded(chatRecord *business.ChatRecord) *business.ChatRecord {
	updatedChatRecord := &business.ChatRecord{}
	if chatRecord != nil {
		*updatedChatRecord = *chatRecord
	}
	return updatedChatRecord
}
