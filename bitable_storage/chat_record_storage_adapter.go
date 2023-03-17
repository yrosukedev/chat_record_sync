package bitable_storage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/yrosukedev/chat_record_sync/business"
	"github.com/yrosukedev/chat_record_sync/consts"
	"net/http"
	"strings"
)

type ChatRecordStorageAdapter struct {
	ctx        context.Context
	larkClient *lark.Client
	appToken   string
	tableId    string
}

func NewChatRecordStorageAdapter(ctx context.Context, larkClient *lark.Client, appToken string, tableId string) *ChatRecordStorageAdapter {
	return &ChatRecordStorageAdapter{
		ctx:        ctx,
		larkClient: larkClient,
		appToken:   appToken,
		tableId:    tableId,
	}
}

// Write store the chat record to Lark Bitable.
//
// requestUUID is a UUID for retrying the idempotent operation.
// It's very useful for handling the errors of writing the chat record to Bitable.
// Because this operation is idempotent when the UUID is provided, so we could retry the operation many times.
func (c *ChatRecordStorageAdapter) Write(chatRecord *business.ChatRecord, requestUUID string) error {

	fields := c.tableFieldsFrom(chatRecord)

	tableRecord := larkbitable.NewAppTableRecordBuilder().
		Fields(fields).
		Build()

	req := larkbitable.
		NewCreateAppTableRecordReqBuilder().
		AppToken(c.appToken).
		TableId(c.tableId).
		ClientToken(requestUUID).
		UserIdType(larkbitable.UserIdTypeUserId).
		AppTableRecord(tableRecord).
		Build()

	resp, err := c.larkClient.Bitable.AppTableRecord.Create(c.ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fails to create lark bitable table record, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	if resp.Code != 0 {
		return fmt.Errorf("fails to create lark bitable table record, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	return nil
}

func (c *ChatRecordStorageAdapter) tableFieldsFrom(chatRecord *business.ChatRecord) map[string]interface{} {
	fields := map[string]interface{}{
		consts.BitableFieldMsgId:   chatRecord.MsgId,
		consts.BitableFieldAction:  chatRecord.Action,
		consts.BitableFieldFrom:    userToTableField(chatRecord.From),
		consts.BitableFieldTo:      usersToTableField(chatRecord.To),
		consts.BitableFieldRoomId:  chatRecord.RoomId,
		consts.BitableFieldMsgTime: chatRecord.MsgTime.UnixMilli(),
		consts.BitableFieldMsgType: chatRecord.MsgType,
		consts.BitableFieldContent: chatRecord.Content,
	}
	return fields
}

func userToTableField(user *business.User) string {
	return fmt.Sprintf("%v(ID:%v)", user.Name, user.UserId)
}

func usersToTableField(users []*business.User) string {
	var userFields []string
	for _, user := range users {
		userFields = append(userFields, userToTableField(user))
	}

	return strings.Join(userFields, ",")
}