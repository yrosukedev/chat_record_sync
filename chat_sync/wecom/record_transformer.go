package wecom

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
)

type RecordTransformer interface {
	Transform(wecomRecord *ChatRecord) (chatRecord *business.ChatRecord, err error)
}
