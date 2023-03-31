package wecom

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
)

type ChatRecordTransformer interface {
	Transform(wecomChatRecord *ChatRecord, userInfo *UserInfo, externalContacts []*ExternalContact) (record *business.ChatRecord, err error)
}
