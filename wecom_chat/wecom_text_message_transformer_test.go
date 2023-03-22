package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"reflect"
	"testing"
	"time"
)

func TestContent(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
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
		Text: &TextMessage{
			Content: "Hello, there!",
		},
	}
	user := &WeComUserInfo{
		UserID: "id_XuJinSheng",
		Name:   "Xu Jin Sheng",
	}
	contacts := []*WeComExternalContact{
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

	// When
	record, err := transformer.Transform(wecomChatRecord, user, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord, record) {
		t.Errorf("records are not matched, expected: %+v, actual: %+v", expectedRecord, record)
		return
	}
}

func TestContent_nilText(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
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
	}
	user := &WeComUserInfo{
		UserID: "id_XuJinSheng",
		Name:   "Xu Jin Sheng",
	}
	contacts := []*WeComExternalContact{
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
		Content: "",
	}

	// When
	record, err := transformer.Transform(wecomChatRecord, user, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord, record) {
		t.Errorf("records are not matched, expected: %+v, actual: %+v", expectedRecord, record)
		return
	}
}

func TestContent_missMatchedMessageType(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
		Seq:    10,
		MsgID:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From:   "id_XuJinSheng",
		ToList: []string{
			"id_icefog",
		},
		RoomID:  "",
		MsgTime: 1547087894783,
		MsgType: "video",
		Text: &TextMessage{
			Content: "Hello, there!",
		},
	}
	user := &WeComUserInfo{
		UserID: "id_XuJinSheng",
		Name:   "Xu Jin Sheng",
	}
	contacts := []*WeComExternalContact{
		{
			ExternalUserID: "id_icefog",
			Name:           "icefog",
		},
	}

	// When
	_, err := transformer.Transform(wecomChatRecord, user, contacts)
	if !reflect.DeepEqual(err, NewTransformerErrorMessageTypeMissMatched("text", "video")) {
		t.Errorf("error should happen here, \nexpected: %v, \nactual: %v\n", NewTransformerErrorMessageTypeMissMatched("text", "video"), err)
		return
	}
}

func TestUserMissmatched_notFound(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
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
		Text: &TextMessage{
			Content: "Hello, there!",
		},
	}
	user := &WeComUserInfo{
		UserID: "id_xiaowang",
		Name:   "xiao wang",
	}
	contacts := []*WeComExternalContact{
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
			Name:   "<unknown>",
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

	// When
	record, err := transformer.Transform(wecomChatRecord, user, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, \nexpected: %+v, \nactual: %+v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord.From, record.From) {
		t.Errorf("users are not matched, \nexpected: %+v, \nactual: %+v", expectedRecord.From, record.From)
		return
	}
}
