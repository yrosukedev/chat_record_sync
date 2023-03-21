package wecom_chat

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"reflect"
	"testing"
)

func TestZeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Eq(uint64(345)), gomock.Eq(uint64(10))).Return(nil, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Any()).Times(0)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Times(0)
	transformer.EXPECT().Transform(gomock.Any()).Times(0)

	// When
	records, outPageToken, err := readerAdapter.Read(paginated_reader.NewPageToken(345), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if len(records) != 0 {
		t.Errorf("records count not matched, expected: %v, actual: %v", 0, len(records))
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(345)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(345), outPageToken)
		return
	}
}
