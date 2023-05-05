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
	if defaultTransformer == nil {
		panic("defaultTransformer can't be nil")
	}

	return &RecordTransformerFactory{
		messageTypeToTransformer: messageTypeToTransformer,
		defaultTransformer:       defaultTransformer,
	}
}

func (r *RecordTransformerFactory) Transform(wecomRecord *wecom.ChatRecord) (chatRecord *business.ChatRecord, err error) {
	transformer := r.transformerFor(wecomRecord)
	return transformer.Transform(wecomRecord, nil)
}

func (r *RecordTransformerFactory) transformerFor(wecomRecord *wecom.ChatRecord) (transformer FieldTransformer) {
	if r.messageTypeToTransformer == nil {
		return r.defaultTransformer
	}

	transformer, ok := r.messageTypeToTransformer[wecomRecord.MsgType]
	if !ok {
		return r.defaultTransformer
	}

	return transformer
}
