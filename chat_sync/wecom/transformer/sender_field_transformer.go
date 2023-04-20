package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type SenderFieldTransformer struct {
	openAPIService wecom.OpenAPIService
}

func NewSenderFieldTransformer(openAPIService wecom.OpenAPIService) *SenderFieldTransformer {
	return &SenderFieldTransformer{
		openAPIService: openAPIService,
	}
}

func (t *SenderFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	user, err := t.openAPIService.GetUserInfoByID(wecomRecord.From)
	if err != nil {
		return nil, err
	}

	updatedChatRecord.From = &business.User{
		UserId: user.UserID,
		Name:   user.Name,
	}

	return updatedChatRecord, nil
}

func (t *SenderFieldTransformer) copyInputIfNeeded(chatRecord *business.ChatRecord) *business.ChatRecord {
	updatedChatRecord := &business.ChatRecord{}
	if chatRecord != nil {
		*updatedChatRecord = *chatRecord
	}
	return updatedChatRecord
}
