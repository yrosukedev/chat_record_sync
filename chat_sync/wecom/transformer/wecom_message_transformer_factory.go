package transformer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type WeComMessageTransformerFactory struct {
	messageTypeToTransformers map[string]wecom.ChatRecordTransformer
	defaultTransformer        wecom.ChatRecordTransformer
}

func NewWeComMessageTransformerFactory(messageTypeToTransformers map[string]wecom.ChatRecordTransformer, defaultTransformer wecom.ChatRecordTransformer) wecom.ChatRecordTransformer {
	return &WeComMessageTransformerFactory{
		messageTypeToTransformers: messageTypeToTransformers,
		defaultTransformer:        defaultTransformer,
	}
}

func (f *WeComMessageTransformerFactory) Transform(wecomChatRecord *wecom.ChatRecord, userInfo *wecom.UserInfo, externalContacts []*wecom.ExternalContact) (record *business.ChatRecord, err error) {
	if wecomChatRecord == nil {
		return nil, nil
	}

	return f.transformerFor(wecomChatRecord).Transform(wecomChatRecord, userInfo, externalContacts)
}

func (f *WeComMessageTransformerFactory) transformerFor(record *wecom.ChatRecord) wecom.ChatRecordTransformer {
	if transformer, ok := f.messageTypeToTransformers[record.MsgType]; ok {
		return transformer
	}

	return f.defaultTransformer
}
