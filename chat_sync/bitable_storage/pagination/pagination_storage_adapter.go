package pagination

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"github.com/yrosukedev/chat_record_sync/consts"
	"github.com/yrosukedev/chat_record_sync/logger"
	"net/http"
	"strconv"
)

type StorageAdapter struct {
	ctx        context.Context
	larkClient *lark.Client
	appToken   string
	tableId    string
	logger     logger.Logger
}

func NewStorageAdapter(ctx context.Context, larkClient *lark.Client, appToken string, tableId string, logger logger.Logger) *StorageAdapter {
	return &StorageAdapter{
		ctx:        ctx,
		larkClient: larkClient,
		appToken:   appToken,
		tableId:    tableId,
		logger:     logger,
	}
}

func (p *StorageAdapter) Get() (pageToken *pagination.PageToken, err error) {

	p.logger.Info(p.ctx, "[pagination storage adapter] will get page token, appToken: %v, tableId: %v", p.appToken, p.tableId)

	req := p.buildRequestOfFetchingLatestRecord()

	resp, err := p.larkClient.Bitable.AppTableRecord.List(p.ctx, req)
	if err := p.checkListRecordsErrors(err, resp); err != nil {
		p.logger.Error(p.ctx, "[pagination storage adapter] fails to get page token, appToken: %v, tableId: %v, error: %v", p.appToken, p.tableId, err)
		return nil, err
	}

	pageToken, err = p.handleResponse(resp)
	if err != nil {
		p.logger.Error(p.ctx, "[pagination storage adapter] fails to get page token, appToken: %v, tableId: %v, error: %v", p.appToken, p.tableId, err)
		return nil, err
	}

	p.logger.Info(p.ctx, "[pagination storage adapter] succeeds to get page token, appToken: %v, tableId: %v, page token: %v", p.appToken, p.tableId, pageToken)

	return pageToken, err
}

func (p *StorageAdapter) handleResponse(resp *larkbitable.ListAppTableRecordResp) (*pagination.PageToken, error) {
	if len(resp.Data.Items) == 0 { // empty
		return nil, nil
	}

	pageTokenField, ok := resp.Data.Items[0].Fields[consts.BitableFieldPaginationPageToken]
	if !ok {
		return nil, fmt.Errorf("bitable field %v dosen't exist", consts.BitableFieldPaginationPageToken)
	}

	return p.pageTokenFrom(pageTokenField)
}

func (p *StorageAdapter) pageTokenFrom(pageTokenField interface{}) (*pagination.PageToken, error) {
	switch pageTokenValue := pageTokenField.(type) {
	case string:
		if pageTokenInt, err := strconv.ParseInt(pageTokenValue, 10, 64); err != nil {
			return nil, fmt.Errorf("fails to convert '%v' field to integer, value: %v", consts.BitableFieldPaginationPageToken, pageTokenValue)
		} else {
			return pagination.NewPageToken(uint64(pageTokenInt)), nil
		}
	case uint64:
		return pagination.NewPageToken(pageTokenValue), nil
	default:
		return nil, fmt.Errorf("unknown type of '%v' field, value: %v", consts.BitableFieldPaginationPageToken, pageTokenField)
	}
}

func (p *StorageAdapter) checkListRecordsErrors(err error, resp *larkbitable.ListAppTableRecordResp) error {
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fails to get records from lark bitable, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	if resp.Code != 0 {
		return fmt.Errorf("fails to get records from lark bitable, statusCode: %v, code: %v, msg: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	return nil
}

func (p *StorageAdapter) buildRequestOfFetchingLatestRecord() *larkbitable.ListAppTableRecordReq {
	req := larkbitable.
		NewListAppTableRecordReqBuilder().
		AppToken(p.appToken).
		TableId(p.tableId).
		Sort(fmt.Sprintf("[\"%v DESC\"]", consts.BitableFieldPaginationDateCreated)).
		FieldNames(fmt.Sprintf("[\"%v\"]", consts.BitableFieldPaginationPageToken)).
		UserIdType(larkbitable.UserIdTypeUserId).
		PageSize(1).
		Build()
	return req
}

func (p *StorageAdapter) Set(pageToken *pagination.PageToken) error {
	p.logger.Info(p.ctx, "[pagination storage adapter] will set page token, appToken: %v, tableId: %v, page token: %v", p.appToken, p.tableId, pageToken)

	if pageToken == nil { // append nothing
		p.logger.Info(p.ctx, "[pagination storage adapter] succeeds to set page token, appToken: %v, tableId: %v, page token: %v", p.appToken, p.tableId, pageToken)
		return nil
	}

	req := p.buildRequestOfAppendingPageToken(pageToken)

	resp, err := p.larkClient.Bitable.AppTableRecord.Create(p.ctx, req)
	if err := p.checkCreateRecordErrors(err, resp); err != nil {
		p.logger.Info(p.ctx, "[pagination storage adapter] fails to set page token, appToken: %v, tableId: %v, page token: %v, error: %v", p.appToken, p.tableId, pageToken, err)
		return err
	}

	p.logger.Info(p.ctx, "[pagination storage adapter] succeeds to set page token, appToken: %v, tableId: %v, page token: %v, record id: %v", p.appToken, p.tableId, pageToken, resp.Data.Record.RecordId)

	return nil
}

func (p *StorageAdapter) buildRequestOfAppendingPageToken(pageToken *pagination.PageToken) *larkbitable.CreateAppTableRecordReq {
	fields := p.tableFieldsFrom(pageToken)

	tableRecord := larkbitable.NewAppTableRecordBuilder().
		Fields(fields).
		Build()

	req := larkbitable.
		NewCreateAppTableRecordReqBuilder().
		AppToken(p.appToken).
		TableId(p.tableId).
		UserIdType(larkbitable.UserIdTypeUserId).
		AppTableRecord(tableRecord).
		Build()
	return req
}

func (p *StorageAdapter) checkCreateRecordErrors(err error, resp *larkbitable.CreateAppTableRecordResp) error {
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

func (p *StorageAdapter) tableFieldsFrom(pageToken *pagination.PageToken) map[string]interface{} {
	return map[string]interface{}{
		consts.BitableFieldPaginationPageToken: fmt.Sprintf("%d", pageToken.Value),
	}
}
