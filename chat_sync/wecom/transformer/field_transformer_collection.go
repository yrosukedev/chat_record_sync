package transformer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type FieldTransformerCollection struct {
	transformers []FieldTransformer
}

func NewFieldTransformerCollection(transformers []FieldTransformer) *FieldTransformerCollection {
	return &FieldTransformerCollection{
		transformers: transformers,
	}
}

func (c *FieldTransformerCollection) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if len(c.transformers) == 0 {
		return chatRecord, nil
	}

	updatedChatRecord = chatRecord
	for _, transformer := range c.transformers {
		tmpChatRecord, err := transformer.Transform(wecomRecord, updatedChatRecord)
		if err == nil {
			updatedChatRecord = tmpChatRecord
		}
	}

	return updatedChatRecord, nil
}
