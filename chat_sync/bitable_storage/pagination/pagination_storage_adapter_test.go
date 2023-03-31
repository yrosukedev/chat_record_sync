//go:build integration
// +build integration

package pagination

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"github.com/yrosukedev/chat_record_sync/config"
	"testing"
)

func TestGetPageToken_succeeds(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	logger := NewMockLogger(ctrl)
	paginationStorage := NewStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3", logger)

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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
	ctrl := gomock.NewController(t)
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	logger := NewMockLogger(ctrl)
	paginationStorage := NewStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3", logger)

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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
	ctrl := gomock.NewController(t)
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	logger := NewMockLogger(ctrl)
	paginationStorage := NewStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3", logger)

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	if err := paginationStorage.Set(pagination.NewPageToken(789478)); err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
	}
}
