package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"io"
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

func TestReceiverFieldTransformer_Transform_wecomRecordCantBeNil(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}

func TestReceiverFieldTransformer_Transform_dontChangeInputs(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		ToList: []string{"123"},
	}
	chatRecord := &business.ChatRecord{}
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
	updatedChatRecord, err := transformer.Transform(wecomRecord, chatRecord)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, updatedChatRecord)

		// Make sure the original chat record is not changed
		assert.Equal(t, chatRecord, &business.ChatRecord{})
	}
}

func TestReceiverFieldTransformer_Transform_zeroReceiver(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{}
	expectedChatRecord := &business.ChatRecord{}

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestReceiverFieldTransformer_Transform_manyReceivers(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		ToList: []string{"123", "456", "789"},
	}
	expectedChatRecord := &business.ChatRecord{
		To: []*business.User{
			{
				UserId: "123",
				Name:   "haary",
			},
			{
				UserId: "456",
				Name:   "lucy",
			},
			{
				UserId: "789",
				Name:   "lily",
			},
		},
	}

	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("123")).Return(&wecom.ExternalContact{
		ExternalUserID: "123",
		Name:           "haary",
	}, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("456")).Return(&wecom.ExternalContact{
		ExternalUserID: "456",
		Name:           "lucy",
	}, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("789")).Return(&wecom.ExternalContact{
		ExternalUserID: "789",
		Name:           "lily",
	}, nil).Times(1)

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestReceiverFieldTransformer_Transform_partialFailure(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewReceiverFieldTransformer(openAPIService)
	wecomRecord := &wecom.ChatRecord{
		ToList: []string{"123", "456", "789"},
	}
	expectedChatRecord := &business.ChatRecord{
		To: []*business.User{
			{
				UserId: "123",
				Name:   "haary",
			},
			{
				UserId: "456",
			},
			{
				UserId: "789",
				Name:   "lily",
			},
		},
	}

	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("123")).Return(&wecom.ExternalContact{
		ExternalUserID: "123",
		Name:           "haary",
	}, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("456")).Return(nil, io.ErrClosedPipe).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("789")).Return(&wecom.ExternalContact{
		ExternalUserID: "789",
		Name:           "lily",
	}, nil).Times(1)

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}
