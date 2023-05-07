package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type SenderFieldTransformer struct {
	nameFetcher NameFetcher
}

func NewSenderFieldTransformer(nameFetcher NameFetcher) *SenderFieldTransformer {
	return &SenderFieldTransformer{
		nameFetcher: nameFetcher,
	}
}

func (t *SenderFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	name, err := t.nameFetcher.FetchName(wecomRecord.From)

	// fatal tolerated
	// logging is done in openAPIService
	if err != nil {
		updatedChatRecord.From = &business.User{
			UserId: wecomRecord.From,
		}
	} else {
		updatedChatRecord.From = &business.User{
			UserId: wecomRecord.From,
			Name:   name,
		}
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
