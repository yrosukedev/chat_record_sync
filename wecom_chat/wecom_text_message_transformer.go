package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"time"
)

type WeComTextMessageTransformer struct {
}

func NewWeComTextMessageTransformer() ChatRecordTransformer {
	return &WeComTextMessageTransformer{}
}

func (w *WeComTextMessageTransformer) Transform(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo, externalContacts []*WeComExternalContact) (record *business.ChatRecord, err error) {
	fromUser := &business.User{
		UserId: userInfo.UserID,
		Name:   userInfo.Name,
	}

	var toUsers []*business.User
	for _, contact := range externalContacts {
		toUsers = append(toUsers, &business.User{
			UserId: contact.ExternalUserID,
			Name:   contact.Name,
		})
	}

	record = &business.ChatRecord{
		MsgId:   wecomChatRecord.MsgID,
		Action:  wecomChatRecord.Action,
		From:    fromUser,
		To:      toUsers,
		RoomId:  wecomChatRecord.RoomID,
		MsgTime: time.UnixMilli(wecomChatRecord.MsgTime),
		MsgType: wecomChatRecord.MsgType,
		Content: wecomChatRecord.Text.Content,
	}

	return record, nil
}
