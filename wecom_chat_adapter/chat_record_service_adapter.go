package wecom_chat_adapter

import (
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
	"time"
)

type WeComChatRecordServiceAdapter struct {
	client  WeWorkFinanceSDK.Client
	proxy   string
	passwd  string
	timeout time.Duration
}

func NewWeComChatRecordServiceAdapter(client WeWorkFinanceSDK.Client, proxy string, passwd string, timeout time.Duration) wecom_chat.ChatRecordService {
	return &WeComChatRecordServiceAdapter{
		client:  client,
		proxy:   proxy,
		passwd:  passwd,
		timeout: timeout,
	}
}

func (w *WeComChatRecordServiceAdapter) Read(seq uint64, pageSize uint64) (records []*wecom_chat.WeComChatRecord, err error) {
	chatDataList, err := w.client.GetChatData(seq, pageSize, w.proxy, w.passwd, int(w.timeout.Seconds()))

	if err != nil {
		return nil, err
	}

	for _, chatData := range chatDataList {
		chatMessage, err := w.client.DecryptData(chatData.EncryptRandomKey, chatData.EncryptChatMsg)
		if err != nil {
			return nil, err
		}

		switch chatMessage.Type {
		case wecom_chat.WeComMessageTypeText:
			records = append(records, w.textMessageFrom(chatMessage, chatData.Seq))
		default:
			records = append(records, w.otherMessageFrom(chatMessage, chatData.Seq))
		}
	}

	return records, nil
}

func (w *WeComChatRecordServiceAdapter) textMessageFrom(chatMessage WeWorkFinanceSDK.ChatMessage, messageSeq uint64) *wecom_chat.WeComChatRecord {
	textMsg := chatMessage.GetTextMessage()
	record := &wecom_chat.WeComChatRecord{
		Seq:     messageSeq,
		MsgID:   textMsg.MsgID,
		Action:  textMsg.Action,
		From:    textMsg.From,
		ToList:  textMsg.ToList,
		RoomID:  textMsg.RoomID,
		MsgTime: textMsg.MsgTime,
		Text: &wecom_chat.TextMessage{
			Content: textMsg.Text.Content,
		},
		OriginMessage: chatMessage.GetRawChatMessage(),
	}
	return record
}

func (w *WeComChatRecordServiceAdapter) otherMessageFrom(chatMessage WeWorkFinanceSDK.ChatMessage, messageSeq uint64) *wecom_chat.WeComChatRecord {
	record := &wecom_chat.WeComChatRecord{
		Seq:           messageSeq,
		MsgID:         chatMessage.Id,
		Action:        chatMessage.Action,
		From:          chatMessage.From,
		ToList:        chatMessage.ToList,
		RoomID:        chatMessage.RoomID,
		MsgTime:       chatMessage.MsgTime,
		OriginMessage: chatMessage.GetRawChatMessage(),
	}
	return record
}
