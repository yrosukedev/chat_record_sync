package wecom_chat

import "github.com/yrosukedev/chat_record_sync/business"

type WeComMessageTransformerFactory struct {
	messageTypeToTransformers map[string]ChatRecordTransformer
}

func NewWeComMessageTransformerFactory(messageTypeToTransformers map[string]ChatRecordTransformer, defaultTransformer ChatRecordTransformer) ChatRecordTransformer {
	return &WeComMessageTransformerFactory{
		messageTypeToTransformers: messageTypeToTransformers,
	}
}

func (f *WeComMessageTransformerFactory) Transform(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo, externalContacts []*WeComExternalContact) (record *business.ChatRecord, err error) {
	return f.transformerFor(wecomChatRecord).Transform(wecomChatRecord, userInfo, externalContacts)
}

func (f *WeComMessageTransformerFactory) transformerFor(record *WeComChatRecord) ChatRecordTransformer {
	return f.messageTypeToTransformers[record.MsgType]
}
