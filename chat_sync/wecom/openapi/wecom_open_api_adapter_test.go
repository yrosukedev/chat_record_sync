//go:build integration
// +build integration

package openapi

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/config"
	"testing"
)

func TestWeComOpenAPIAdapter_GetUserInfoByID_success(t *testing.T) {
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
	if assert.NoError(t, err) {
		assert.Equal(t, expectedUserInfo, userInfo)
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
	assert.Error(t, err)
	assert.Nil(t, userInfo)

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
	if assert.NoError(t, err) {
		assert.Equal(t, expectedExternalContact, externalContact)
	}
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
	assert.Error(t, err)
	assert.Nil(t, externalContact)

	// Then
}

// test external room success case
func TestWeComOpenAPIAdapter_GetExternalRoomByID_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.AgentSecret, wecomConfig.AgentID)
	logger := NewMockLogger(ctrl)
	openAPI := NewAdapter(ctx, wecomApp, logger)

	roomId := "wrsrtEBgAANfNIS5R-b8uWKPkS3S0Y6w"
	expectedExternalRoom := &wecom.ExternalRoom{
		RoomID: roomId,
		Name:   "会话存档测试群",
	}

	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// When
	room, err := openAPI.GetExternalRoomByID(roomId)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedExternalRoom, room)
	}
}
