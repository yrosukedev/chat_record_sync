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
	if wecomChatRecord.MsgType != WeComMessageTypeText {
		return nil, NewTransformerErrorMessageTypeMissMatched(WeComMessageTypeText, wecomChatRecord.MsgType)
	}

	fromUser := w.userFrom(wecomChatRecord, userInfo)

	var toUsers []*business.User
	for _, contact := range externalContacts {
		toUsers = append(toUsers, &business.User{
			UserId: contact.ExternalUserID,
			Name:   contact.Name,
		})
	}

	content := w.contentFrom(wecomChatRecord)

	record = &business.ChatRecord{
		MsgId:   wecomChatRecord.MsgID,
		Action:  wecomChatRecord.Action,
		From:    fromUser,
		To:      toUsers,
		RoomId:  wecomChatRecord.RoomID,
		MsgTime: time.UnixMilli(wecomChatRecord.MsgTime),
		MsgType: wecomChatRecord.MsgType,
		Content: content,
	}

	return record, nil
}

func (w *WeComTextMessageTransformer) userFrom(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo) *business.User {
	userName := "<unknown>"
	if wecomChatRecord.From == userInfo.UserID {
		userName = userInfo.Name
	}

	fromUser := &business.User{
		UserId: wecomChatRecord.From,
		Name:   userName,
	}
	return fromUser
}

func (w *WeComTextMessageTransformer) contentFrom(wecomChatRecord *WeComChatRecord) string {
	result := ""
	if wecomChatRecord.Text != nil {
		result = wecomChatRecord.Text.Content
	}
	return result
}
