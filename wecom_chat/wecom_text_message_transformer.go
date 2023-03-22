package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"time"
)

type WeComMessageTransformer struct {
	contentTransformer func(wecomChatRecord *WeComChatRecord) (content string, err error)
}

func NewWeComTextMessageTransformer() ChatRecordTransformer {
	return &WeComMessageTransformer{
		contentTransformer: func(wecomChatRecord *WeComChatRecord) (content string, err error) {
			if wecomChatRecord.MsgType != WeComMessageTypeText {
				return "", NewTransformerErrorMessageTypeMissMatched(WeComMessageTypeText, wecomChatRecord.MsgType)
			}

			content = ""
			if wecomChatRecord.Text != nil {
				content = wecomChatRecord.Text.Content
			}
			return content, nil
		},
	}
}

func (w *WeComMessageTransformer) Transform(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo, externalContacts []*WeComExternalContact) (record *business.ChatRecord, err error) {
	if wecomChatRecord == nil {
		return nil, nil
	}

	content, err := w.contentTransformer(wecomChatRecord)
	if err != nil {
		return nil, err
	}

	sender := w.senderFrom(wecomChatRecord, userInfo)

	receivers := w.receiversFrom(wecomChatRecord, externalContacts)

	record = &business.ChatRecord{
		MsgId:   wecomChatRecord.MsgID,
		Action:  wecomChatRecord.Action,
		From:    sender,
		To:      receivers,
		RoomId:  wecomChatRecord.RoomID,
		MsgTime: time.UnixMilli(wecomChatRecord.MsgTime),
		MsgType: wecomChatRecord.MsgType,
		Content: content,
	}

	return record, nil
}

func (w *WeComMessageTransformer) receiversFrom(wecomChatRecord *WeComChatRecord, externalContacts []*WeComExternalContact) []*business.User {
	contactIdToNames := make(map[string]string, len(externalContacts))
	for _, contact := range externalContacts {
		contactIdToNames[contact.ExternalUserID] = contact.Name
	}

	var results []*business.User
	for _, contactId := range wecomChatRecord.ToList {
		contactName := "<unknown>"
		if name, ok := contactIdToNames[contactId]; ok {
			contactName = name
		}

		results = append(results, &business.User{
			UserId: contactId,
			Name:   contactName,
		})
	}
	return results
}

func (w *WeComMessageTransformer) senderFrom(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo) *business.User {
	userName := w.senderNameFrom(wecomChatRecord, userInfo)

	fromUser := &business.User{
		UserId: wecomChatRecord.From,
		Name:   userName,
	}
	return fromUser
}

func (w *WeComMessageTransformer) senderNameFrom(wecomChatRecord *WeComChatRecord, userInfo *WeComUserInfo) string {
	result := "<unknown>"
	if userInfo != nil && userInfo.UserID == wecomChatRecord.From {
		result = userInfo.Name
	}
	return result
}
