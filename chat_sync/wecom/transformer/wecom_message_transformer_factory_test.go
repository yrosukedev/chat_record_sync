package transformer

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"reflect"
	"testing"
	"time"
)

func TestTransformerFactory_transformerFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	transformer1 := wecom.NewMockChatRecordTransformer(ctrl)
	transformerForTextMsg := wecom.NewMockChatRecordTransformer(ctrl)
	transformer3 := wecom.NewMockChatRecordTransformer(ctrl)
	messageTypeToTransformers := map[string]wecom.ChatRecordTransformer{
		"::any type but text 1::": transformer1,
		wecom.MessageTypeText:     transformerForTextMsg,
		"::any type but text 2::": transformer3,
	}
	factory := NewWeComMessageTransformerFactory(messageTypeToTransformers, nil)
	wecomChatRecord, user, contacts, expectedRecord := makeSUTForTransformerFactory()

	// Then
	transformer1.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	transformer3.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	transformerForTextMsg.EXPECT().Transform(gomock.Eq(wecomChatRecord), gomock.Eq(user), gomock.Eq(contacts)).Return(expectedRecord, nil).Times(1)

	// When
	record, err := factory.Transform(wecomChatRecord, user, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord, record) {
		t.Errorf("records are not matched, expected: %+v, actual: %+v", expectedRecord, record)
		return
	}
}

func TestTransformerFactory_transformerNotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	transformer1 := wecom.NewMockChatRecordTransformer(ctrl)
	transformer2 := wecom.NewMockChatRecordTransformer(ctrl)
	transformer3 := wecom.NewMockChatRecordTransformer(ctrl)
	defaultTransformer := wecom.NewMockChatRecordTransformer(ctrl)
	messageTypeToTransformers := map[string]wecom.ChatRecordTransformer{
		"::any type but text 1::": transformer1,
		"::any type but text 2::": transformer2,
		"::any type but text 3::": transformer3,
	}
	factory := NewWeComMessageTransformerFactory(messageTypeToTransformers, defaultTransformer)
	wecomChatRecord, user, contacts, expectedRecord := makeSUTForTransformerFactory()

	// Then
	transformer1.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	transformer2.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	transformer3.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	defaultTransformer.EXPECT().Transform(gomock.Eq(wecomChatRecord), gomock.Eq(user), gomock.Eq(contacts)).Return(expectedRecord, nil).Times(1)

	// When
	record, err := factory.Transform(wecomChatRecord, user, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord, record) {
		t.Errorf("records are not matched, expected: %+v, actual: %+v", expectedRecord, record)
		return
	}
}

func TestTransformerFactory_nilParameters(t *testing.T) {
	// Given
	factory := NewWeComMessageTransformerFactory(nil, nil)

	// When
	record, err := factory.Transform(nil, nil, nil)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if nil != record {
		t.Errorf("records are not matched, expected: %+v, actual: %+v", nil, record)
		return
	}
}

func makeSUTForTransformerFactory() (*wecom.ChatRecord, *wecom.UserInfo, []*wecom.ExternalContact, *business.ChatRecord) {
	wecomChatRecord := &wecom.ChatRecord{
		Seq:    10,
		MsgID:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From:   "id_XuJinSheng",
		ToList: []string{
			"id_icefog",
		},
		RoomID:  "",
		MsgTime: 1547087894783,
		MsgType: "text",
		Text: &wecom.TextMessage{
			Content: "Hello, there!",
		},
	}
	user := &wecom.UserInfo{
		UserID: "id_XuJinSheng",
		Name:   "Xu Jin Sheng",
	}
	contacts := []*wecom.ExternalContact{
		{
			ExternalUserID: "id_icefog",
			Name:           "icefog",
		},
	}
	expectedRecord := &business.ChatRecord{
		MsgId:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From: &business.User{
			UserId: "id_XuJinSheng",
			Name:   "Xu Jin Sheng",
		},
		To: []*business.User{
			{
				UserId: "id_icefog",
				Name:   "icefog",
			},
		},
		RoomId:  "",
		MsgTime: time.UnixMilli(1547087894783),
		MsgType: "text",
		Content: "Hello, there!",
	}
	return wecomChatRecord, user, contacts, expectedRecord
}
