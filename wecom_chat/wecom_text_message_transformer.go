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

	sender := w.senderFrom(wecomChatRecord, userInfo)

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
		From:    sender,
		To:      toUsers,
		RoomId:  wecomChatRecord.RoomID,
		MsgTime: time.UnixMilli(wecomChatRecord.MsgTime),
		MsgType: wecomChatRecord.MsgType,
		Content: content,
	}

	return record, nil
}

func (w *WeComTextMessageTransformer) senderFrom(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo) *business.User {
	userName := w.senderNameFrom(wecomChatRecord, userInfo)

	fromUser := &business.User{
		UserId: wecomChatRecord.From,
		Name:   userName,
	}
	return fromUser
}

func (w *WeComTextMessageTransformer) senderNameFrom(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo) string {
	result := "<unknown>"
	if userInfo != nil && userInfo.UserID == wecomChatRecord.From {
		result = userInfo.Name
	}
	return result
}

func (w *WeComTextMessageTransformer) contentFrom(wecomChatRecord *WeComChatRecord) string {
	result := ""
	if wecomChatRecord.Text != nil {
		result = wecomChatRecord.Text.Content
	}
	return result
}
