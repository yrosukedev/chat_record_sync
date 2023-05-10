package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type OtherContentFieldTransformer struct {
}

func NewOtherContentFieldTransformer() *OtherContentFieldTransformer {
	return &OtherContentFieldTransformer{}
}

func (t *OtherContentFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	if wecomRecord.OriginMessage == nil {
		return nil, wecom.NewTransformerEmptyContentError(wecomRecord)
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	updatedChatRecord.Content = string(wecomRecord.OriginMessage)

	return updatedChatRecord, nil
}

func (t *OtherContentFieldTransformer) copyInputIfNeeded(chatRecord *business.ChatRecord) *business.ChatRecord {
	updatedChatRecord := &business.ChatRecord{}
	if chatRecord != nil {
		*updatedChatRecord = *chatRecord
	}
	return updatedChatRecord
}
