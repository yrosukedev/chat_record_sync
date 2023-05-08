package transformer

import (
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
)

type RoomFieldTransformer struct {
	nameFetcher NameFetcher
}

func NewRoomFieldTransformer(nameFetcher NameFetcher) *RoomFieldTransformer {
	return &RoomFieldTransformer{
		nameFetcher: nameFetcher,
	}
}

func (t *RoomFieldTransformer) Transform(wecomRecord *wecom.ChatRecord, chatRecord *business.ChatRecord) (updatedChatRecord *business.ChatRecord, err error) {
	if wecomRecord == nil {
		return nil, errors.New("wecomRecord can't be nil")
	}

	updatedChatRecord = t.copyInputIfNeeded(chatRecord)

	name, err := t.nameFetcher.FetchName(wecomRecord.RoomID)

	// fatal tolerated
	// logging is done in openAPIService
	if err != nil {
		updatedChatRecord.Room = &business.Room{
			RoomId: wecomRecord.RoomID,
		}
	} else {
		updatedChatRecord.Room = &business.Room{
			RoomId: wecomRecord.RoomID,
			Name:   name,
		}
	}

	return updatedChatRecord, nil
}

func (t *RoomFieldTransformer) copyInputIfNeeded(chatRecord *business.ChatRecord) *business.ChatRecord {
	updatedChatRecord := &business.ChatRecord{}
	if chatRecord != nil {
		*updatedChatRecord = *chatRecord
	}
	return updatedChatRecord
}
