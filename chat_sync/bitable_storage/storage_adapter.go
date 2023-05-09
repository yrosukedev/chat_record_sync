package bitable_storage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/yrosukedev/chat_record_sync/logger"
	"net/http"
)

type StorageAdapter struct {
	ctx        context.Context
	larkClient *lark.Client
	appToken   string
	tableId    string
	logger     logger.Logger
}

func NewStorageAdapter(
	ctx context.Context,
	larkClient *lark.Client,
	appToken string,
	tableId string,
	logger logger.Logger) *StorageAdapter {
	return &StorageAdapter{
		ctx:        ctx,
		larkClient: larkClient,
		appToken:   appToken,
		tableId:    tableId,
		logger:     logger,
	}
}

func (s *StorageAdapter) Write(fields map[string]interface{}, requestUUID string) error {
	tableRecord := larkbitable.NewAppTableRecordBuilder().
		Fields(fields).
		Build()

	req := larkbitable.
		NewCreateAppTableRecordReqBuilder().
		AppToken(s.appToken).
		TableId(s.tableId).
		ClientToken(requestUUID).
		UserIdType(larkbitable.UserIdTypeUserId).
		AppTableRecord(tableRecord).
		Build()

	resp, err := s.larkClient.Bitable.AppTableRecord.Create(s.ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fails to create lark bitable table record, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
		return err
	}

	if resp.Code != 0 {
		err = fmt.Errorf("fails to create lark bitable table record, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
		return err
	}

	return nil
}



