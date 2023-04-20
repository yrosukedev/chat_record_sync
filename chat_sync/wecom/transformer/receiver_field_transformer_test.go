package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestReceiverFieldTransformer_Transform_nilChatRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		ToList: []string{"123"},
	}
	expectedChatRecord := &business.ChatRecord{
		To: []*business.User{
			{
				UserId: "123",
				Name:   "haary",
			},
		},
	}

	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("123")).Return(&wecom.ExternalContact{
		ExternalUserID: "123",
		Name:           "haary",
	}, nil).Times(1)

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
