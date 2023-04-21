package transformer

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"testing"
)

func TestOtherContentFieldTransformer_Transform_nilChatRecord(t *testing.T) {
	// Given
	transformer := NewOtherContentFieldTransformer()
	originMessage := "{\"msgid\":\"2641513858500683770_1603876152\",\"action\":\"send\",\"from\":\"icefog\",\"tolist\":[\"wmN6etBgAA0sbJ3invMvRxPQDFoq9uWA\"],\"roomid\":\"\",\"msgtime\":1603876152141,\"msgtype\":\"location\",\"location\":{\"longitude\":116.586285899,\"latitude\":39.911125799,\"address\":\"北京市xxx区xxx路xxx大厦x座\",\"title\":\"xxx管理中心\",\"zoom\":15}}"
	wecomRecord := &wecom.ChatRecord{
		MsgType:       "location",
		OriginMessage: []byte(originMessage),
	}
	expectedChatRecord := &business.ChatRecord{
		Content: originMessage,
	}

	// When
	chatRecord, err := transformer.Transform(wecomRecord, nil)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedChatRecord, chatRecord)
	}
}

func TestOtherContentFieldTransformer_Transform_wecomRecordCantBeNil(t *testing.T) {
	// Given
	transformer := NewOtherContentFieldTransformer()

	// When
	chatRecord, err := transformer.Transform(nil, &business.ChatRecord{})

	// Then
	if assert.Error(t, err) {
		assert.Nil(t, chatRecord)
	}
}
