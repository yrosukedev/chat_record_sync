package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"reflect"
	"testing"
	"time"
)

func TestTextMessage_Content(t *testing.T) {
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

func TestOtherMessage_Content(t *testing.T) {
	// Given
	transformer := NewWeComDefaultMessageTransformer()
	originMessage := "{\"msgid\":\"2641513858500683770_1603876152\",\"action\":\"send\",\"from\":\"icefog\",\"tolist\":[\"wmN6etBgAA0sbJ3invMvRxPQDFoq9uWA\"],\"roomid\":\"\",\"msgtime\":1603876152141,\"msgtype\":\"location\",\"location\":{\"longitude\":116.586285899,\"latitude\":39.911125799,\"address\":\"北京市xxx区xxx路xxx大厦x座\",\"title\":\"xxx管理中心\",\"zoom\":15}}"
	wecomChatRecord := &WeComChatRecord{
		Seq:    10,
		MsgID:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From:   "id_XuJinSheng",
		ToList: []string{
			"id_icefog",
		},
		RoomID:        "",
		MsgTime:       1547087894783,
		MsgType:       "location",
		OriginMessage: []byte(originMessage),
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
		MsgType: "location",
		Content: originMessage,
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

func TestTextMessage_Content_nilText(t *testing.T) {
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

func TestOtherMessage_Content_nilOriginMessage(t *testing.T) {
	// Given
	transformer := NewWeComDefaultMessageTransformer()
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
		MsgType: "location",
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
	if !reflect.DeepEqual(err, NewTransformerEmptyContentError(wecomChatRecord)) {
		t.Errorf("error should happen here, expected: %v, actual: %v", NewTransformerEmptyContentError(wecomChatRecord), err)
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

func TestUserMissmatched_nilUser(t *testing.T) {
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
	record, err := transformer.Transform(wecomChatRecord, nil, contacts)
	if err != nil {
		t.Errorf("error shouldn't happen here, \nexpected: %+v, \nactual: %+v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord.From, record.From) {
		t.Errorf("users are not matched, \nexpected: %+v, \nactual: %+v", expectedRecord.From, record.From)
		return
	}
}

func TestContactsMissmatched_contactNotFound(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
		Seq:    10,
		MsgID:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From:   "id_XuJinSheng",
		ToList: []string{
			"id_icefog",
			"id_xiaohuang",
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
			ExternalUserID: "id_xiaohuang",
			Name:           "Xiao Huang",
		},
		{
			ExternalUserID: "id_xiaozhang",
			Name:           "Xiao Zhang",
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
				Name:   "<unknown>",
			},
			{
				UserId: "id_xiaohuang",
				Name:   "Xiao Huang",
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

	if !reflect.DeepEqual(expectedRecord.To, record.To) {
		t.Errorf("contacts are not matched, \nexpected: %+v, \nactual: %+v", expectedRecord.To, record.To)
		return
	}
}

func TestContactsMissmatched_nilContact(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()
	wecomChatRecord := &WeComChatRecord{
		Seq:    10,
		MsgID:  "CAQQluDa4QUY0On2rYSAgAMgzPrShAE=",
		Action: "send",
		From:   "id_XuJinSheng",
		ToList: []string{
			"id_icefog",
			"id_xiaohuang",
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
				Name:   "<unknown>",
			},
			{
				UserId: "id_xiaohuang",
				Name:   "<unknown>",
			},
		},
		RoomId:  "",
		MsgTime: time.UnixMilli(1547087894783),
		MsgType: "text",
		Content: "Hello, there!",
	}

	// When
	record, err := transformer.Transform(wecomChatRecord, user, nil)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if !reflect.DeepEqual(expectedRecord.To, record.To) {
		t.Errorf("contacts are not matched, \nexpected: %+v, \nactual: %+v", expectedRecord.To, record.To)
		return
	}
}

func TestNilWeComChatRecord(t *testing.T) {
	// Given
	transformer := NewWeComTextMessageTransformer()

	// When
	record, err := transformer.Transform(nil, nil, nil)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	if record != nil {
		t.Errorf("records not matched, expected: %+v, actual: %+v", nil, record)
		return
	}
}