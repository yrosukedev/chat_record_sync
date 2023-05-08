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

// test internal room success case
func TestMsgAuditOpenAPIAdapter_GetInternalRoomByID_Success(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	wecomConfig := config.NewWeComConfig()
	wecomApp := workwx.New(wecomConfig.CorpID).WithApp(wecomConfig.ChatSyncSecret, config.WeComMsgAuditAgentID)
	logger := NewMockLogger(ctrl)
	openAPI := NewMsgAuditOpenAPIAdapter(ctx, wecomApp, logger)

	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	roomId := "wrsrtEBgAAp36zvfx_uMyjEyBkXEmyMQ"
	expectedRoom := &wecom.InternalRoom{
		RoomID: roomId,
		Name:   "内部群测试~",
	}

	// When
	room, err := openAPI.GetInternalRoomByID(roomId)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedRoom, room)
	}
}
