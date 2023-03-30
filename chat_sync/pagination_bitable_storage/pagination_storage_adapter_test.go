//go:build integration
// +build integration

package pagination_bitable_storage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/yrosukedev/chat_record_sync/chat_sync/paginated_reader"
	"github.com/yrosukedev/chat_record_sync/config"
	"testing"
)

func TestGetPageToken_succeeds(t *testing.T) {
	// Given
	ctx := context.Background()
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	paginationStorage := NewPaginationStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3")

	// When
	token, err := paginationStorage.Get()

	// Then
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %+v, actual: %+v", nil, err)
	}

	fmt.Printf("page token: %v\n", token)
}

func TestGetPageToken_nil(t *testing.T) {
	// Given
	ctx := context.Background()
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	paginationStorage := NewPaginationStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3")

	// When
	token, err := paginationStorage.Get()

	// Then
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %+v, actual: %+v", nil, err)
	}

	if token != nil {
		t.Errorf("page token is expected to be %v, actual: %v", nil, token)
	}
}

func TestSetPageToken_succeeds(t *testing.T) {
	// Given
	ctx := context.Background()
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	paginationStorage := NewPaginationStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3")

	// When
	if err := paginationStorage.Set(paginated_reader.NewPageToken(789478)); err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
	}
}
