//go:build integration
// +build integration

package openapi

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

	openAPI := NewAdapter(ctx, wecomApp, logger)
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

// test failure case
func TestWeComOpenAPIAdapter_GetUserInfoByID_Failure(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)

	logger := NewMockLogger(ctrl)

	openAPI := NewAdapter(ctx, wecomApp, logger)
	userId := "user ID that doesn't exist"

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	userInfo, err := openAPI.GetUserInfoByID(userId)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if userInfo != nil {
		t.Errorf("user info should be nil, expected: %#v, actual: %#v", nil, userInfo)
		return
	}

	// Then
}

// test external contact success case
func TestWeComOpenAPIAdapter_GetExternalContactByID_success(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)

	logger := NewMockLogger(ctrl)

	openAPI := NewAdapter(ctx, wecomApp, logger)
	externalId := "wmsrtEBgAAJEj0gFLCHAfuv75QLoLmgw"
	expectedExternalContact := &wecom.ExternalContact{
		ExternalUserID: externalId,
		Name:           "王欢",
	}

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	externalContact, err := openAPI.GetExternalContactByID(externalId)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedExternalContact, externalContact) {
		t.Errorf("external contact are not matched, expected: %#v, actual: %#v", expectedExternalContact, externalContact)
		return
	}

	// Then
}

// test external contact failure case
func TestWeComOpenAPIAdapter_GetExternalContactByID_Failure(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)

	logger := NewMockLogger(ctrl)

	openAPI := NewAdapter(ctx, wecomApp, logger)
	externalId := "external ID that doesn't exist"

	// When
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	externalContact, err := openAPI.GetExternalContactByID(externalId)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if externalContact != nil {
		t.Errorf("external contact should be nil, expected: %#v, actual: %#v", nil, externalContact)
		return
	}

	// Then
}
