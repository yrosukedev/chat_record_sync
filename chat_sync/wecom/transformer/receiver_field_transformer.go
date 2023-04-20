package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type ReceiverFieldTransformer struct {
	openAPIService wecom.OpenAPIService
}

func NewReceiverFieldTransformer(openAPIService wecom.OpenAPIService) *ReceiverFieldTransformer {
	return &ReceiverFieldTransformer{
		openAPIService: openAPIService,
	}
}

func (t *ReceiverFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	// transform each of the receivers in wecomRecord.ToList to business.User by calling openAPIService
	// and append the result to updatedChatRecord.To
	// if any error occurs, ignore it and continue
	for _, receiverId := range wecomRecord.ToList {
		contact, err := t.openAPIService.GetExternalContactByID(receiverId)
		if err == nil {
			updatedChatRecord.To = append(updatedChatRecord.To, &business.User{
				UserId: contact.ExternalUserID,
				Name:   contact.Name,
			})
		}
	}

	return updatedChatRecord, nil
}

func (t *ReceiverFieldTransformer) copyInputIfNeeded(chatRecord *business.ChatRecord) *business.ChatRecord {
	updatedChatRecord := &business.ChatRecord{}
	if chatRecord != nil {
		*updatedChatRecord = *chatRecord
	}
	return updatedChatRecord
}
