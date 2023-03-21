package wecom_chat

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
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
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

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

func TestOneRecord_oneReceiver(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	wecomRecords := []*WeComChatRecord{
		{
			Seq:    890,
			From:   "ID_xiaoming",
			ToList: []string{"ID_xiaowang"},
		},
	}
	user := &WeComUserInfo{
		UserID: "ID_xiaoming",
		Name:   "Xiao Ming",
	}
	contacts := []*WeComExternalContact{
		{
			ExternalUserID: "ID_xiaowang",
			Name:           "Xiao Wang",
		},
	}
	expectedRecords := []*business.ChatRecord{
		{
			From: &business.User{
				UserId: "ID_xiaoming",
				Name:   "Xiao Ming",
			},
			To: []*business.User{
				{
					UserId: "ID_xiaowang",
					Name:   "Xiao Wang",
				},
			},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(user, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaowang")).Return(contacts[0], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[0]), gomock.Eq(user), gomock.Eq(contacts)).Return(expectedRecords[0], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(paginated_reader.NewPageToken(267), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(890)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(890), outPageToken)
		return
	}
}

func TestOneRecord_zeroReceiver(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	wecomRecords := []*WeComChatRecord{
		{
			Seq:  375,
			From: "ID_xiaoming",
		},
	}
	user := &WeComUserInfo{
		UserID: "ID_xiaoming",
		Name:   "Xiao Ming",
	}
	expectedRecords := []*business.ChatRecord{
		{
			From: &business.User{
				UserId: "ID_xiaoming",
				Name:   "Xiao Ming",
			},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(user, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Times(0)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[0]), gomock.Eq(user), gomock.Nil()).Return(expectedRecords[0], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(paginated_reader.NewPageToken(267), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(375)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(375), outPageToken)
		return
	}
}
