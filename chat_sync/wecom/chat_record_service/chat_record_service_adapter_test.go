//go:build integration
// +build integration

package chat_record_service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/config"
	"testing"
)

func TestWeComChatRecordServiceSDK_firstPage(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	weComConfig := config.NewWeComConfig()
	client, err := WeWorkFinanceSDK.NewClient(weComConfig.CorpID, weComConfig.ChatSyncSecret, weComConfig.ChatSyncRsaPrivateKey)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	logger := NewMockLogger(ctrl)

	adapter := NewAdapter(ctx, client, "", "", config.WeComChatRecordSDKTimeout, logger)

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	records, err := adapter.Read(0, 10)

	// Then
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	fmt.Println("WeCom Records:")
	for idx, record := range records {
		fmt.Printf("[%v] %#v\n", idx, record)
	}
}
