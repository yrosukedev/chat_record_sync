package pagination_bitable_storage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	"github.com/yrosukedev/chat_record_sync/consts"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"io"
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

	req := p.buildRequestByFetchingLatestRecord()

	resp, err := p.larkClient.Bitable.AppTableRecord.List(p.ctx, req)

	if err := p.checkErrors(err, resp); err != nil {
		return paginated_reader.NewPageToken(0), err
	}

	return p.handleResponse(resp)
}

func (p *PaginationStorageAdapter) handleResponse(resp *larkbitable.ListAppTableRecordResp) (*paginated_reader.PageToken, error) {
	if len(resp.Data.Items) == 0 { // empty
		return paginated_reader.NewPageToken(0), nil
	}

	pageTokenField, ok := resp.Data.Items[0].Fields[consts.BitableFieldPaginationPageToken]
	if !ok {
		return paginated_reader.NewPageToken(0), fmt.Errorf("bitable field %v dosen't exist, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, p.appToken, p.tableId)
	}

	return p.pageTokenFrom(pageTokenField)
}

func (p *PaginationStorageAdapter) pageTokenFrom(pageTokenField interface{}) (*paginated_reader.PageToken, error) {
	switch pageTokenValue := pageTokenField.(type) {
	case string:
		if pageTokenInt, err := strconv.Atoi(pageTokenValue); err != nil {
			return paginated_reader.NewPageToken(0), fmt.Errorf("fails to convert '%v' field to integer, value: %v, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, pageTokenValue, p.appToken, p.tableId)
		} else {
			return paginated_reader.NewPageToken(int64(pageTokenInt)), nil
		}
	case int:
		return paginated_reader.NewPageToken(int64(pageTokenValue)), nil
	default:
		return paginated_reader.NewPageToken(0), fmt.Errorf("unknown type of '%v' field, value: %v, appToken: %v, tableId: %v", consts.BitableFieldPaginationPageToken, pageTokenField, p.appToken, p.tableId)
	}
}

func (p *PaginationStorageAdapter) checkErrors(err error, resp *larkbitable.ListAppTableRecordResp) error {
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

func (p *PaginationStorageAdapter) buildRequestByFetchingLatestRecord() *larkbitable.ListAppTableRecordReq {
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
	// TODO:
	return io.ErrClosedPipe
}
