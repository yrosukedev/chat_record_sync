package pagination_bitable_storage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/yrosukedev/chat_record_sync/consts"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"net/http"
	"strconv"
)

type PaginationStorageAdapter struct {
	ctx        context.Context
	larkClient *lark.Client
	appToken   string
	tableId    string
}

func NewPaginationStorageAdapter(ctx context.Context, larkClient *lark.Client, appToken string, tableId string) *PaginationStorageAdapter {
	return &PaginationStorageAdapter{
		ctx:        ctx,
		larkClient: larkClient,
		appToken:   appToken,
		tableId:    tableId,
	}
}

func (p *PaginationStorageAdapter) Get() (pageToken *paginated_reader.PageToken, err error) {

	req := p.buildRequestOfFetchingLatestRecord()

	resp, err := p.larkClient.Bitable.AppTableRecord.List(p.ctx, req)

	if err := p.checkListRecordsErrors(err, resp); err != nil {
		return nil, err
	}

	return p.handleResponse(resp)
}

func (p *PaginationStorageAdapter) handleResponse(resp *larkbitable.ListAppTableRecordResp) (*paginated_reader.PageToken, error) {
	if len(resp.Data.Items) == 0 { // empty
		return nil, nil
	}

	pageTokenField, ok := resp.Data.Items[0].Fields[consts.BitableFieldPaginationPageToken]
	if !ok {
		return nil, fmt.Errorf("bitable field %v dosen't exist, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, p.appToken, p.tableId)
	}

	return p.pageTokenFrom(pageTokenField)
}

func (p *PaginationStorageAdapter) pageTokenFrom(pageTokenField interface{}) (*paginated_reader.PageToken, error) {
	switch pageTokenValue := pageTokenField.(type) {
	case string:
		if pageTokenInt, err := strconv.Atoi(pageTokenValue); err != nil {
			return nil, fmt.Errorf("fails to convert '%v' field to integer, value: %v, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, pageTokenValue, p.appToken, p.tableId)
		} else {
			return paginated_reader.NewPageToken(int64(pageTokenInt)), nil
		}
	case int:
		return paginated_reader.NewPageToken(int64(pageTokenValue)), nil
	default:
		return nil, fmt.Errorf("unknown type of '%v' field, value: %v, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, pageTokenField, p.appToken, p.tableId)
	}
}

func (p *PaginationStorageAdapter) checkListRecordsErrors(err error, resp *larkbitable.ListAppTableRecordResp) error {
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

func (p *PaginationStorageAdapter) buildRequestOfFetchingLatestRecord() *larkbitable.ListAppTableRecordReq {
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

func (p *PaginationStorageAdapter) Set(pageToken *paginated_reader.PageToken) error {
	if pageToken == nil { // append nothing
		return nil
	}

	req := p.buildRequestOfAppendingPageToken(pageToken)

	resp, err := p.larkClient.Bitable.AppTableRecord.Create(p.ctx, req)

	if err := p.checkCreateRecordErrors(err, resp); err != nil {
		return err
	}

	return nil
}

func (p *PaginationStorageAdapter) buildRequestOfAppendingPageToken(pageToken *paginated_reader.PageToken) *larkbitable.CreateAppTableRecordReq {
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

func (p *PaginationStorageAdapter) checkCreateRecordErrors(err error, resp *larkbitable.CreateAppTableRecordResp) error {
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

func (p *PaginationStorageAdapter) tableFieldsFrom(pageToken *paginated_reader.PageToken) map[string]interface{} {
	return map[string]interface{}{
		consts.BitableFieldPaginationPageToken: fmt.Sprintf("%d", pageToken.Value),
	}
}
