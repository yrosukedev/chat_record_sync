//go:build integration
// +build integration

package wecom_openapi_adapter

import (
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/config"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
	"reflect"
	"testing"
)

func TestWeComOpenAPIAdapter_GetUserInfoByID(t *testing.T) {
	// Given
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)

	openAPI := NewWeComOpenAPIAdapter(wecomApp)
	userId := "WangHuan"
	expectedUserInfo := &wecom_chat.WeComUserInfo{
		UserID: userId,
		Name:   "王欢",
	}

	// When
	userInfo, err := openAPI.GetUserInfoByID(userId)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedUserInfo, userInfo) {
		t.Errorf("user info are not matched, expected: %#v, actual: %#v", expectedUserInfo, userInfo)
		return
	}

	// Then
}
