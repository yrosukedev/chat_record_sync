package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestSenderFieldTransformer_Transform_nilChatRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewSenderFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		From: "123",
	}
	expectedChatRecord := &business.ChatRecord{
		From: &business.User{
			UserId: "123",
			Name:   "haary",
		},
	}

	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("123")).Return(&wecom.UserInfo{
		UserID: "123",
		Name:   "haary",
	}, nil).Times(1)

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
