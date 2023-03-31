package transformer

import (
	"context"
	"fmt"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/logger"
	"strings"
	"time"
)

type WeComMessageTransformer struct {
	contentTransformer func(wecomChatRecord *wecom.ChatRecord) (content string, err error)
	logger             logger.Logger
	ctx                context.Context
}

func NewWeComTextMessageTransformer(ctx context.Context, logger logger.Logger) wecom.ChatRecordTransformer {
	return &WeComMessageTransformer{
		contentTransformer: func(wecomChatRecord *wecom.ChatRecord) (content string, err error) {
			if wecomChatRecord.MsgType != wecom.MessageTypeText {
				panic(fmt.Sprintf("missmatched message type, expected: %v, actual: %v", wecom.MessageTypeText, wecomChatRecord.MsgType))
			}

			content = ""
			if wecomChatRecord.Text != nil {
				content = wecomChatRecord.Text.Content
			}
			return content, nil
		},
		logger: logger,
		ctx:    ctx,
	}
}

func NewWeComDefaultMessageTransformer(ctx context.Context, logger logger.Logger) wecom.ChatRecordTransformer {
	return &WeComMessageTransformer{
		contentTransformer: func(wecomChatRecord *wecom.ChatRecord) (content string, err error) {
			content = ""
			if wecomChatRecord.OriginMessage != nil {
				content = string(wecomChatRecord.OriginMessage)
			}

			if content == "" {
				return "", wecom.NewTransformerEmptyContentError(wecomChatRecord)
			}

			return content, nil
		},
		logger: logger,
		ctx:    ctx,
	}
}

func (w *WeComMessageTransformer) Transform(wecomChatRecord *wecom.ChatRecord, userInfo *wecom.UserInfo, externalContacts []*wecom.ExternalContact) (record *business.ChatRecord, err error) {
	if wecomChatRecord == nil {
		w.logger.Info(w.ctx, "[message transformer] message is nil, nothing to do")
		return nil, nil
	}

	w.logger.Info(w.ctx, "[message transformer] will transform message, seq: %v, msgId: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID)

	content, err := w.contentTransformer(wecomChatRecord)
	if err != nil {
		w.logger.Error(w.ctx, "[message transformer] fails to transform message content, seq: %v, msgId: %v, error: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID, err)
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

	w.logger.Info(w.ctx, "[message transformer] did transform message, seq: %v, msgId: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID)

	return record, nil
}

func (w *WeComMessageTransformer) receiversFrom(wecomChatRecord *wecom.ChatRecord, externalContacts []*wecom.ExternalContact) []*business.User {
	contactIdToNames := make(map[string]string, len(externalContacts))
	for _, contact := range externalContacts {
		contactIdToNames[contact.ExternalUserID] = contact.Name
	}

	var results []*business.User
	var missingContactIds []string
	for _, contactId := range wecomChatRecord.ToList {
		name, ok := contactIdToNames[contactId]
		if !ok {
			missingContactIds = append(missingContactIds, contactId)
			name = "<unknown>"
		}

		results = append(results, &business.User{
			UserId: contactId,
			Name:   name,
		})
	}

	if len(missingContactIds) > 0 {
		w.logger.Error(w.ctx, "[message transformer] contact ids are not found when transforming contact id to contact name, seq: %v, msgId: %v, missing contact ids: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID, strings.Join(missingContactIds, ","))
	}

	return results
}

func (w *WeComMessageTransformer) senderFrom(wecomChatRecord *wecom.ChatRecord, userInfo *wecom.UserInfo) *business.User {
	userName := w.senderNameFrom(wecomChatRecord, userInfo)

	fromUser := &business.User{
		UserId: wecomChatRecord.From,
		Name:   userName,
	}
	return fromUser
}

func (w *WeComMessageTransformer) senderNameFrom(wecomChatRecord *wecom.ChatRecord, userInfo *wecom.UserInfo) string {
	if userInfo == nil {
		w.logger.Error(w.ctx, "[message transformer] fails to transform user id to user name, seq: %v, msgId: %v, error: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID, "user info is nil")
		return "<unknown>"
	}

	if userInfo.UserID != wecomChatRecord.From {
		w.logger.Error(w.ctx, "[message transformer] user ids are not matched when transforming user id to user name, seq: %v, msgId: %v, expected: %v, actual: %v", wecomChatRecord.Seq, wecomChatRecord.MsgID, wecomChatRecord.From, userInfo.UserID)
		return "<unknown>"
	}

	return userInfo.Name
}
