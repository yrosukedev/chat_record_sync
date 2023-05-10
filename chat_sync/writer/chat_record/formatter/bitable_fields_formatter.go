package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/consts"
	"golang.org/x/exp/maps"
	"strings"
)

type FieldsCreator = func(*business.ChatRecord) (map[string]interface{}, error)

type BitableFieldsFormatter struct {
}

func NewBitableFieldsFormatter() *BitableFieldsFormatter {
	return &BitableFieldsFormatter{}
}

func (b *BitableFieldsFormatter) Format(record *business.ChatRecord) (fields map[string]interface{}, err error) {
	fields = make(map[string]interface{})

	for _, creator := range b.fieldCreators() {
		partialFields, err := creator(record)
		if err != nil {
			return nil, err
		}
		maps.Copy(fields, partialFields)
	}

	return fields, nil
}

func (b *BitableFieldsFormatter) fieldCreators() []FieldsCreator {
	fieldsCreator := []FieldsCreator{
		b.basicFields,
		b.senderFields,
		b.receiverFields,
		b.roomFields,
	}
	return fieldsCreator
}

func (b *BitableFieldsFormatter) roomFields(record *business.ChatRecord) (map[string]interface{}, error) {
	if record.Room == nil {
		return map[string]interface{}{
			consts.BitableFieldChatRecordRoomId:   "",
			consts.BitableFieldChatRecordRoomName: "",
		}, nil
	}

	return map[string]interface{}{
		consts.BitableFieldChatRecordRoomId:   record.Room.RoomId,
		consts.BitableFieldChatRecordRoomName: record.Room.Name,
	}, nil
}

func (b *BitableFieldsFormatter) receiverFields(record *business.ChatRecord) (fields map[string]interface{}, err error) {
	var receiverIds []string
	receiverIdToName := make(map[string]string)

	for _, user := range record.To {
		receiverIds = append(receiverIds, user.UserId)
		receiverIdToName[user.UserId] = user.Name
	}

	receiverNamesJson, err := json.Marshal(receiverIdToName)
	if err != nil {
		return nil, fmt.Errorf("fails to marshal receiver names, receiverIdToName: %v, err: %v", receiverIdToName, err)
	}

	return map[string]interface{}{
		consts.BitableFieldChatRecordReceiverIds:   strings.Join(receiverIds, ","),
		consts.BitableFieldChatRecordReceiverNames: string(receiverNamesJson),
	}, nil
}

func (b *BitableFieldsFormatter) senderFields(record *business.ChatRecord) (map[string]interface{}, error) {
	if record.From == nil {
		return map[string]interface{}{
			consts.BitableFieldChatRecordSenderId:   "",
			consts.BitableFieldChatRecordSenderName: "",
		}, nil
	}

	return map[string]interface{}{
		consts.BitableFieldChatRecordSenderId:   record.From.UserId,
		consts.BitableFieldChatRecordSenderName: record.From.Name,
	}, nil
}

func (b *BitableFieldsFormatter) basicFields(record *business.ChatRecord) (map[string]interface{}, error) {
	return map[string]interface{}{
		consts.BitableFieldChatRecordMsgId:   record.MsgId,
		consts.BitableFieldChatRecordAction:  record.Action,
		consts.BitableFieldChatRecordMsgTime: record.MsgTime.UnixMilli(),
		consts.BitableFieldChatRecordMsgType: record.MsgType,
		consts.BitableFieldChatRecordContent: record.Content,
	}, nil
}
