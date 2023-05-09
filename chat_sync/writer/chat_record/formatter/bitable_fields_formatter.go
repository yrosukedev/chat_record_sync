package formatter

import (
	"fmt"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/consts"
	"golang.org/x/exp/maps"
	"strings"
)

type FieldsCreator = func(*business.ChatRecord) map[string]interface{}

type BitableFieldsFormatter struct {
}

func NewBitableFieldsFormatter() *BitableFieldsFormatter {
	return &BitableFieldsFormatter{}
}

func (b *BitableFieldsFormatter) Format(record *business.ChatRecord) (fields map[string]interface{}, err error) {
	fields = make(map[string]interface{})

	for _, creator := range b.fieldCreators() {
		maps.Copy(fields, creator(record))
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

func (b *BitableFieldsFormatter) roomFields(record *business.ChatRecord) map[string]interface{} {
	if record.Room == nil {
		return map[string]interface{}{
			consts.BitableFieldChatRecordRoomId:   "",
			consts.BitableFieldChatRecordRoomName: "",
		}
	}

	return map[string]interface{}{
		consts.BitableFieldChatRecordRoomId:   record.Room.RoomId,
		consts.BitableFieldChatRecordRoomName: record.Room.Name,
	}
}

func (b *BitableFieldsFormatter) receiverFields(record *business.ChatRecord) map[string]interface{} {
	return map[string]interface{}{
		consts.BitableFieldChatRecordTo: b.usersToField(record.To),
	}
}

func (b *BitableFieldsFormatter) senderFields(record *business.ChatRecord) map[string]interface{} {
	return map[string]interface{}{
		consts.BitableFieldChatRecordFrom: b.userToField(record.From),
	}
}

func (b *BitableFieldsFormatter) basicFields(record *business.ChatRecord) map[string]interface{} {
	return map[string]interface{}{
		consts.BitableFieldChatRecordMsgId:   record.MsgId,
		consts.BitableFieldChatRecordAction:  record.Action,
		consts.BitableFieldChatRecordMsgTime: record.MsgTime.UnixMilli(),
		consts.BitableFieldChatRecordMsgType: record.MsgType,
		consts.BitableFieldChatRecordContent: record.Content,
	}
}

func (b *BitableFieldsFormatter) userToField(user *business.User) string {
	if user == nil {
		return ""
	}
	return fmt.Sprintf("%v(ID:%v)", user.Name, user.UserId)
}

func (b *BitableFieldsFormatter) usersToField(users []*business.User) string {
	var userFields []string
	for _, user := range users {
		userFields = append(userFields, b.userToField(user))
	}

	return strings.Join(userFields, ",")
}
