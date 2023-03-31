//go:build integration
// +build integration

package wecom_openapi_adapter

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/config"
	"reflect"
	"testing"
)

func TestWeComOpenAPIAdapter_GetUserInfoByID(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)

	logger := NewMockLogger(ctrl)

	openAPI := NewWeComOpenAPIAdapter(ctx, wecomApp, logger)
	userId := "WangHuan"
	expectedUserInfo := &wecom.UserInfo{
		UserID: userId,
		Name:   "王欢",
	}

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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
