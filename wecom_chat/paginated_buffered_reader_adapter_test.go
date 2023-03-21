package wecom_chat

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"io"
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

func TestZeroRecord_nilInputPageToken(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Eq(uint64(0)), gomock.Eq(uint64(10))).Return(nil, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Any()).Times(0)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Times(0)
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// When
	records, outPageToken, err := readerAdapter.Read(nil, 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if len(records) != 0 {
		t.Errorf("records count not matched, expected: %v, actual: %v", 0, len(records))
		return
	}
	if outPageToken != nil {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", nil, outPageToken)
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

func TestOneRecord_manyReceivers(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	wecomRecords := []*WeComChatRecord{
		{
			Seq:  1334,
			From: "ID_xiaoming",
			ToList: []string{
				"ID_xiaowang",
				"ID_xiaozhang",
				"ID_xiaoli",
			},
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
		{
			ExternalUserID: "ID_xiaozhang",
			Name:           "Xiao Zhang",
		},
		{
			ExternalUserID: "ID_xiaoli",
			Name:           "Xiao Li",
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
				{
					UserId: "ID_xiaozhang",
					Name:   "Xiao Zhang",
				},
				{
					UserId: "ID_xiaoli",
					Name:   "Xiao Li",
				},
			},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(123)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(user, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaowang")).Return(contacts[0], nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaozhang")).Return(contacts[1], nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaoli")).Return(contacts[2], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[0]), gomock.Eq(user), gomock.Eq(contacts)).Return(expectedRecords[0], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(paginated_reader.NewPageToken(123), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(1334)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(1334), outPageToken)
		return
	}
}

func TestManyRecords(t *testing.T) {
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
		{
			Seq:    357,
			From:   "ID_xiaohuang",
			ToList: []string{"ID_xiaoli"},
		},
		{
			Seq:    1050,
			From:   "ID_xiaoming",
			ToList: []string{"ID_xiaozhang"},
		},
	}
	userXiaoming := &WeComUserInfo{
		UserID: "ID_xiaoming",
		Name:   "Xiao Ming",
	}
	userXiaohuang := &WeComUserInfo{
		UserID: "ID_xiaohuang",
		Name:   "Xiao Huang",
	}
	contactXiaowang := &WeComExternalContact{
		ExternalUserID: "ID_xiaowang",
		Name:           "Xiao Wang",
	}
	contactXiaoli := &WeComExternalContact{
		ExternalUserID: "ID_xiaoli",
		Name:           "Xiao Li",
	}
	contactXiaozhang := &WeComExternalContact{
		ExternalUserID: "ID_xiaozhang",
		Name:           "Xiao Zhang",
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
		{
			From: &business.User{
				UserId: "ID_xiaohuang",
				Name:   "Xiao Huang",
			},
			To: []*business.User{
				{
					UserId: "ID_xiaoli",
					Name:   "Xiao Li",
				},
			},
		},
		{
			From: &business.User{
				UserId: "ID_xiaoming",
				Name:   "Xiao Ming",
			},
			To: []*business.User{
				{
					UserId: "ID_xiaozhang",
					Name:   "Xiao Zhang",
				},
			},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(15)), gomock.Eq(uint64(30))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(userXiaoming, nil).Times(2)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaohuang")).Return(userXiaohuang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaowang")).Return(contactXiaowang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaozhang")).Return(contactXiaozhang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaoli")).Return(contactXiaoli, nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[0]), gomock.Eq(userXiaoming), gomock.Eq([]*WeComExternalContact{contactXiaowang})).Return(expectedRecords[0], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[1]), gomock.Eq(userXiaohuang), gomock.Eq([]*WeComExternalContact{contactXiaoli})).Return(expectedRecords[1], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[2]), gomock.Eq(userXiaoming), gomock.Eq([]*WeComExternalContact{contactXiaozhang})).Return(expectedRecords[2], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(paginated_reader.NewPageToken(15), 30)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(1050)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(1050), outPageToken)
		return
	}
}

func TestManyRecords_nilInputPageToken(t *testing.T) {
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
		{
			Seq:    357,
			From:   "ID_xiaohuang",
			ToList: []string{"ID_xiaoli"},
		},
		{
			Seq:    1050,
			From:   "ID_xiaoming",
			ToList: []string{"ID_xiaozhang"},
		},
	}
	userXiaoming := &WeComUserInfo{
		UserID: "ID_xiaoming",
		Name:   "Xiao Ming",
	}
	userXiaohuang := &WeComUserInfo{
		UserID: "ID_xiaohuang",
		Name:   "Xiao Huang",
	}
	contactXiaowang := &WeComExternalContact{
		ExternalUserID: "ID_xiaowang",
		Name:           "Xiao Wang",
	}
	contactXiaoli := &WeComExternalContact{
		ExternalUserID: "ID_xiaoli",
		Name:           "Xiao Li",
	}
	contactXiaozhang := &WeComExternalContact{
		ExternalUserID: "ID_xiaozhang",
		Name:           "Xiao Zhang",
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
		{
			From: &business.User{
				UserId: "ID_xiaohuang",
				Name:   "Xiao Huang",
			},
			To: []*business.User{
				{
					UserId: "ID_xiaoli",
					Name:   "Xiao Li",
				},
			},
		},
		{
			From: &business.User{
				UserId: "ID_xiaoming",
				Name:   "Xiao Ming",
			},
			To: []*business.User{
				{
					UserId: "ID_xiaozhang",
					Name:   "Xiao Zhang",
				},
			},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(0)), gomock.Eq(uint64(30))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(userXiaoming, nil).Times(2)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaohuang")).Return(userXiaohuang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaowang")).Return(contactXiaowang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaozhang")).Return(contactXiaozhang, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaoli")).Return(contactXiaoli, nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[0]), gomock.Eq(userXiaoming), gomock.Eq([]*WeComExternalContact{contactXiaowang})).Return(expectedRecords[0], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[1]), gomock.Eq(userXiaohuang), gomock.Eq([]*WeComExternalContact{contactXiaoli})).Return(expectedRecords[1], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Eq(wecomRecords[2]), gomock.Eq(userXiaoming), gomock.Eq([]*WeComExternalContact{contactXiaozhang})).Return(expectedRecords[2], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(nil, 30)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, paginated_reader.NewPageToken(1050)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", paginated_reader.NewPageToken(1050), outPageToken)
		return
	}
}

func TestChatRecordServiceError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	openAPIService := NewMockOpenAPIService(ctrl)
	transformer := NewMockChatRecordTransformer(ctrl)
	readerAdapter := NewPaginatedBufferedReaderAdapter(chatRecordService, openAPIService, transformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil, io.ErrShortBuffer).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Any()).Times(0)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Times(0)
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// When
	_, _, err := readerAdapter.Read(paginated_reader.NewPageToken(345), 10)
	if err == nil {
		t.Errorf("error shouldn happen here, expected: %v, actual: %v", io.ErrShortBuffer, err)
		return
	}
}

func TestOpenAPIServiceError_getUserInfo(t *testing.T) {
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

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(nil, io.ErrClosedPipe).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Times(0)
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// When
	_, _, err := readerAdapter.Read(paginated_reader.NewPageToken(267), 10)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", io.ErrClosedPipe, err)
		return
	}
}

func TestOpenAPIServiceError_getContact(t *testing.T) {
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

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(user, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Any()).Return(nil, io.ErrNoProgress).Times(1)
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// When
	_, _, err := readerAdapter.Read(paginated_reader.NewPageToken(267), 10)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", io.ErrNoProgress, err)
		return
	}
}

func TestTransformerError(t *testing.T) {
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

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	openAPIService.EXPECT().GetUserInfoByID(gomock.Eq("ID_xiaoming")).Return(user, nil).Times(1)
	openAPIService.EXPECT().GetExternalContactByID(gomock.Eq("ID_xiaowang")).Return(contacts[0], nil).Times(1)
	transformer.EXPECT().Transform(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, io.ErrUnexpectedEOF).Times(1)

	// When
	_, _, err := readerAdapter.Read(paginated_reader.NewPageToken(267), 10)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", io.ErrUnexpectedEOF, err)
		return
	}
}
