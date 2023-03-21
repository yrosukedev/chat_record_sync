package wecom_chat

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordTransformer interface {
	Transform(wecomChatRecord *WeComChatRecord) (record *business.ChatRecord, err error)
}
