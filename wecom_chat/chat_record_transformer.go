package wecom_chat

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordTransformer interface {
	Transform(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo, externalContacts []*WeComExternalContact) (record *business.ChatRecord, err error)
}
