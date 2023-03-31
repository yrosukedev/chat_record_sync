package wecom_chat_adapter

import (
	"context"
	"fmt"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/logger"
	"strings"
	"time"
)

type WeComChatRecordServiceAdapter struct {
	ctx     context.Context
	client  WeWorkFinanceSDK.Client
	proxy   string
	passwd  string
	timeout time.Duration
	logger  logger.Logger
}

func NewWeComChatRecordServiceAdapter(
	ctx context.Context,
	client WeWorkFinanceSDK.Client,
	proxy string,
	passwd string,
	timeout time.Duration,
	logger logger.Logger) wecom_chat2.ChatRecordService {
	return &WeComChatRecordServiceAdapter{
		ctx:     ctx,
		client:  client,
		proxy:   proxy,
		passwd:  passwd,
		timeout: timeout,
		logger:  logger,
	}
}

func (w *WeComChatRecordServiceAdapter) Read(seq uint64, pageSize uint64) (records []*wecom_chat2.ChatRecord, err error) {
	w.logger.Info(w.ctx, "[wecom chat record service adapter] will read records, seq: %v, page size: %v", seq, pageSize)
	chatDataList, err := w.client.GetChatData(seq, pageSize, w.proxy, w.passwd, int(w.timeout.Seconds()))
	if err != nil {
		w.logger.Error(w.ctx, "[wecom chat record service adapter] fails to read records, seq: %v, page size: %v, error: %v", seq, pageSize, err)
		return nil, err
	}
	w.logger.Info(w.ctx, "[wecom chat record service adapter] succeeds to read records, seq: %v, page size: %v, chat data list: %v", seq, pageSize, w.summarizeChatDataList(chatDataList))

	w.logger.Info(w.ctx, "[wecom chat record service adapter] will decrypt messages")
	chatDataAndMessages, err := w.decryptMessagesFrom(chatDataList)
	if err != nil {
		return nil, err
	}
	w.logger.Info(w.ctx, "[wecom chat record service adapter] succeeds to decrypt messages")

	w.logger.Info(w.ctx, "[wecom chat record service adapter] will transform messages")
	records = w.transformChatRecordsFrom(chatDataAndMessages)
	w.logger.Info(w.ctx, "[wecom chat record service adapter] succeeds to transform messages")

	return records, nil
}

func (w *WeComChatRecordServiceAdapter) transformChatRecordsFrom(chatDataAndMessages []chatDataAndMessage) []*wecom_chat2.ChatRecord {
	var records []*wecom_chat2.ChatRecord
	for _, chatDataAndMessage := range chatDataAndMessages {
		switch chatDataAndMessage.ChatMessage.Type {
		case wecom_chat2.MessageTypeText:
			records = append(records, w.textMessageFrom(chatDataAndMessage.ChatMessage, chatDataAndMessage.ChatData.Seq))
		default:
			records = append(records, w.otherMessageFrom(chatDataAndMessage.ChatMessage, chatDataAndMessage.ChatData.Seq))
		}
	}
	return records
}

func (w *WeComChatRecordServiceAdapter) decryptMessagesFrom(chatDataList []WeWorkFinanceSDK.ChatData) ([]chatDataAndMessage, error) {
	var chatDataAndMessages []chatDataAndMessage
	for _, chatData := range chatDataList {
		chatMessage, err := w.client.DecryptData(chatData.EncryptRandomKey, chatData.EncryptChatMsg)
		if err != nil {
			w.logger.Error(w.ctx, "[wecom chat record service adapter] fails to decrypt message, seq: %v, msgId: %v, encrypt random key: %v, encrypt msg: %v, error: %v", chatData.Seq, chatData.MsgId, chatData.EncryptRandomKey, chatData.EncryptChatMsg, err)
			return nil, err
		}

		chatDataAndMessages = append(chatDataAndMessages, chatDataAndMessage{chatData, chatMessage})
	}
	return chatDataAndMessages, nil
}

func (w *WeComChatRecordServiceAdapter) summarizeChatDataList(chatDataList []WeWorkFinanceSDK.ChatData) string {
	var result []string
	for _, chatData := range chatDataList {
		result = append(result, fmt.Sprintf("{seq: %v, msgId: %v}", chatData.Seq, chatData.MsgId))
	}
	return strings.Join(result, ",")
}

func (w *WeComChatRecordServiceAdapter) textMessageFrom(chatMessage WeWorkFinanceSDK.ChatMessage, messageSeq uint64) *wecom_chat2.ChatRecord {
	textMsg := chatMessage.GetTextMessage()
	record := &wecom_chat2.ChatRecord{
		Seq:     messageSeq,
		MsgID:   textMsg.MsgID,
		Action:  textMsg.Action,
		From:    textMsg.From,
		ToList:  textMsg.ToList,
		RoomID:  textMsg.RoomID,
		MsgTime: textMsg.MsgTime,
		MsgType: textMsg.MsgType,
		Text: &wecom_chat2.TextMessage{
			Content: textMsg.Text.Content,
		},
		OriginMessage: chatMessage.GetRawChatMessage(),
	}
	return record
}

func (w *WeComChatRecordServiceAdapter) otherMessageFrom(chatMessage WeWorkFinanceSDK.ChatMessage, messageSeq uint64) *wecom_chat2.ChatRecord {
	record := &wecom_chat2.ChatRecord{
		Seq:           messageSeq,
		MsgID:         chatMessage.Id,
		Action:        chatMessage.Action,
		From:          chatMessage.From,
		ToList:        chatMessage.ToList,
		RoomID:        chatMessage.RoomID,
		MsgTime:       chatMessage.MsgTime,
		MsgType:       chatMessage.Type,
		OriginMessage: chatMessage.GetRawChatMessage(),
	}
	return record
}

type chatDataAndMessage struct {
	ChatData    WeWorkFinanceSDK.ChatData
	ChatMessage WeWorkFinanceSDK.ChatMessage
}
