package formatter

import (
	"fmt"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/consts"
	"strings"
)

type BitableFieldsFormatter struct {
}

func NewBitableFieldsFormatter() *BitableFieldsFormatter {
	return &BitableFieldsFormatter{}
}

func (b *BitableFieldsFormatter) Format(record *business.ChatRecord) (fields map[string]interface{}, err error) {
	fields = map[string]interface{}{
		consts.BitableFieldChatRecordMsgId:   record.MsgId,
		consts.BitableFieldChatRecordAction:  record.Action,
		consts.BitableFieldChatRecordFrom:    b.userToTableField(record.From),
		consts.BitableFieldChatRecordTo:      b.usersToTableField(record.To),
		consts.BitableFieldChatRecordRoomId:  record.Room.RoomId,
		consts.BitableFieldChatRecordMsgTime: record.MsgTime.UnixMilli(),
		consts.BitableFieldChatRecordMsgType: record.MsgType,
		consts.BitableFieldChatRecordContent: record.Content,
	}
	return fields, nil
}

func (b *BitableFieldsFormatter) userToTableField(user *business.User) string {
	if user == nil {
		return ""
	}
	return fmt.Sprintf("%v(ID:%v)", user.Name, user.UserId)
}

func (b *BitableFieldsFormatter) usersToTableField(users []*business.User) string {
	var userFields []string
	for _, user := range users {
		userFields = append(userFields, b.userToTableField(user))
	}

	return strings.Join(userFields, ",")
}
