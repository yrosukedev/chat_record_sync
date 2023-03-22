package wecom_chat

import "github.com/yrosukedev/chat_record_sync/business"

type WeComMessageTransformerFactory struct {
	messageTypeToTransformers map[string]ChatRecordTransformer
	defaultTransformer        ChatRecordTransformer
}

func NewWeComMessageTransformerFactory(messageTypeToTransformers map[string]ChatRecordTransformer, defaultTransformer ChatRecordTransformer) ChatRecordTransformer {
	return &WeComMessageTransformerFactory{
		messageTypeToTransformers: messageTypeToTransformers,
		defaultTransformer:        defaultTransformer,
	}
}

func (f *WeComMessageTransformerFactory) Transform(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo, externalContacts []*WeComExternalContact) (record *business.ChatRecord, err error) {
	if wecomChatRecord == nil {
		return nil, nil
	}

	return f.transformerFor(wecomChatRecord).Transform(wecomChatRecord, userInfo, externalContacts)
}

func (f *WeComMessageTransformerFactory) transformerFor(record *WeComChatRecord) ChatRecordTransformer {
	if transformer, ok := f.messageTypeToTransformers[record.MsgType]; ok {
		return transformer
	}

	return f.defaultTransformer
}
