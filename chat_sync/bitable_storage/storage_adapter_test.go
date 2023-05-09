package bitable_storage

import (
	"context"
	"github.com/golang/mock/gomock"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/config"
	"github.com/yrosukedev/chat_record_sync/consts"
	"testing"
)

func TestStorageAdapter_Write(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	logger := NewMockLogger(ctrl)
	storageAdapter := NewStorageAdapter(ctx, larkClient, "QCBrbzgx4aKRAis9eewcV731n7d", "tblIk692K5LXte8x", logger)

	fields := map[string]interface{}{
		consts.BitableFieldChatRecordMsgId:   "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		consts.BitableFieldChatRecordAction:  "send",
		consts.BitableFieldChatRecordFrom:    "Harry Wang(xxyzzwksksk=)",
		consts.BitableFieldChatRecordTo:      "小明(poijnfhdwp=),小黄(mngwscfgyttt=)",
		consts.BitableFieldChatRecordRoomId:  "wmErxtDgAA9AW32YyyuYRimKr7D1KWlw",
		consts.BitableFieldChatRecordMsgTime: int64(1677721600000),
		consts.BitableFieldChatRecordMsgType: "text",
		consts.BitableFieldChatRecordContent: "Let's go take a dinner.",
	}

	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	// When
	err := storageAdapter.Write(fields, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.NoError(t, err)
}
