package transformer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type RecordTransformerFactory struct {
	messageTypeToTransformer map[string]FieldTransformer
	defaultTransformer       FieldTransformer
}

func NewRecordTransformerFactory(messageTypeToTransformer map[string]FieldTransformer, defaultTransformer FieldTransformer) *RecordTransformerFactory {
	return &RecordTransformerFactory{
		messageTypeToTransformer: messageTypeToTransformer,
		defaultTransformer:       defaultTransformer,
	}
}

func (r *RecordTransformerFactory) Transform(wecomRecord *wecom.ChatRecord) (chatRecord *business.ChatRecord, err error) {
	if r.messageTypeToTransformer == nil {
		return nil, nil
	}

	transformer, ok := r.messageTypeToTransformer[wecomRecord.MsgType]
	if !ok {
		transformer = r.defaultTransformer
	}

	return transformer.Transform(wecomRecord, nil)
}
