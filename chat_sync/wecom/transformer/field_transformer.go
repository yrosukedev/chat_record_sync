package transformer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type FieldTransformer interface {
	Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error)
}
