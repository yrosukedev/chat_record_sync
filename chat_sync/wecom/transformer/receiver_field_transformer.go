package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type ReceiverFieldTransformer struct {
	openAPIService OpenAPIService
}

func NewReceiverFieldTransformer(openAPIService OpenAPIService) *ReceiverFieldTransformer {
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
	// if any error occurs, the contact's name fall back to empty value and continue
	for _, receiver := range wecomRecord.ToList {
		contact, err := t.openAPIService.GetExternalContactByID(receiver)

		var user *business.User
		if err != nil {
			// fatal tolerated
			// logging is done in openAPIService
			user = &business.User{
				UserId: receiver,
			}
		} else {
			user = &business.User{
				UserId: contact.ExternalUserID,
				Name:   contact.Name,
			}
		}

		updatedChatRecord.To = append(updatedChatRecord.To, user)
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
