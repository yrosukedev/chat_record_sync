package wecom_chat_adapter

import (
	"fmt"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/config"
	"testing"
)

func TestWeComChatRecordServiceSDK_firstPage(t *testing.T) {
	// Given
	weComConfig := config.NewWeComConfig()
	client, err := WeWorkFinanceSDK.NewClient(weComConfig.CorpID, weComConfig.CorpSecret, weComConfig.RsaPrivateKey)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	adapter := NewWeComChatRecordServiceAdapter(client, "", "", config.WeComChatRecordSDKTimeout)

	// When
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
