package transformer

import "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"

type RecordTransformerBuilder struct {
	openAPIService wecom.OpenAPIService
}

func NewRecordTransformerBuilder(openAPIService wecom.OpenAPIService) *RecordTransformerBuilder {
	return &RecordTransformerBuilder{
		openAPIService: openAPIService,
	}
}

func (b *RecordTransformerBuilder) Build() wecom.RecordTransformer {
	defaultTransformer := NewFieldTransformerCollection([]FieldTransformer{
		NewBasicFieldTransformer(),
		NewSenderFieldTransformer(b.openAPIService),
		NewReceiverFieldTransformer(b.openAPIService),
		NewOtherContentFieldTransformer(),
	})
	textMessageTransformer := NewFieldTransformerCollection([]FieldTransformer{
		NewBasicFieldTransformer(),
		NewSenderFieldTransformer(b.openAPIService),
		NewReceiverFieldTransformer(b.openAPIService),
		NewTextContentFieldTransformer(),
	})
	messageTypeToTransformer := map[string]FieldTransformer{
		wecom.MessageTypeText: textMessageTransformer,
	}
	return NewRecordTransformerFactory(messageTypeToTransformer, defaultTransformer)
}
