package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"io"
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

func TestSenderFieldTransformer_Transform_wecomRecordCantBeNil(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewSenderFieldTransformer(openAPIService)

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}

func TestSenderFieldTransformer_Transform_dontChangeInputs(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewSenderFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		From: "123",
	}
	chatRecord := &business.ChatRecord{
		From: &business.User{
			UserId: "::whatever user id::",
			Name:   "::whatever user name::",
		},
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
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)

		assert.Equal(t, chatRecord.From.UserId, "::whatever user id::")
		assert.Equal(t, chatRecord.From.Name, "::whatever user name::")
	}
}

func TestSenderFieldTransformer_Transform_TolerateError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewSenderFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		From: "123",
	}
	chatRecord := &business.ChatRecord{
		From: &business.User{
			UserId: "::whatever user id::",
			Name:   "::whatever user name::",
		},
	}

	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("123")).Return(nil, io.ErrShortBuffer).Times(1)

	// When
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, updatedChatRecord, chatRecord)
	}
}
