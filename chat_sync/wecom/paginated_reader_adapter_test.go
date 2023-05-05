package wecom

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"io"
	"reflect"
	"testing"
)

func TestZeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Eq(uint64(345)), gomock.Eq(uint64(10))).Return(nil, nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Any()).Times(0)

	// When
	records, outPageToken, err := readerAdapter.Read(pagination.NewPageToken(345), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if len(records) != 0 {
		t.Errorf("records count not matched, expected: %v, actual: %v", 0, len(records))
		return
	}
	if !reflect.DeepEqual(outPageToken, pagination.NewPageToken(345)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", pagination.NewPageToken(345), outPageToken)
		return
	}
}

func TestZeroRecord_nilInputPageToken(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Eq(uint64(0)), gomock.Eq(uint64(10))).Return(nil, nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Any()).Times(0)

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
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	wecomRecords := []*ChatRecord{
		{
			Seq:    890,
			From:   "ID_xiaoming",
			ToList: []string{"ID_xiaowang"},
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
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[0])).Return(expectedRecords[0], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(pagination.NewPageToken(267), 10)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, pagination.NewPageToken(890)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", pagination.NewPageToken(890), outPageToken)
		return
	}
}

func TestManyRecords(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	wecomRecords := []*ChatRecord{
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
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[0])).Return(expectedRecords[0], nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[1])).Return(expectedRecords[1], nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[2])).Return(expectedRecords[2], nil).Times(1)

	// When
	records, outPageToken, err := readerAdapter.Read(pagination.NewPageToken(15), 30)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}
	if !reflect.DeepEqual(expectedRecords, records) {
		t.Errorf("records not matched, expected: %+v, actual: %+v", expectedRecords, records)
		return
	}
	if !reflect.DeepEqual(outPageToken, pagination.NewPageToken(1050)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", pagination.NewPageToken(1050), outPageToken)
		return
	}
}

func TestManyRecords_nilInputPageToken(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	wecomRecords := []*ChatRecord{
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
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[0])).Return(expectedRecords[0], nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[1])).Return(expectedRecords[1], nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Eq(wecomRecords[2])).Return(expectedRecords[2], nil).Times(1)

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
	if !reflect.DeepEqual(outPageToken, pagination.NewPageToken(1050)) {
		t.Errorf("output page token not matched, expected: %+v, actual: %+v", pagination.NewPageToken(1050), outPageToken)
		return
	}
}

func TestChatRecordServiceError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	chatRecordService.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil, io.ErrShortBuffer).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Any()).Times(0)

	// When
	_, _, err := readerAdapter.Read(pagination.NewPageToken(345), 10)
	if err == nil {
		t.Errorf("error shouldn happen here, expected: %v, actual: %v", io.ErrShortBuffer, err)
		return
	}
}

func TestTransformerError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	chatRecordService := NewMockChatRecordService(ctrl)
	recordTransformer := NewMockRecordTransformer(ctrl)
	readerAdapter := NewPaginatedReaderAdapter(chatRecordService, recordTransformer)

	// Then
	wecomRecords := []*ChatRecord{
		{
			Seq:    890,
			From:   "ID_xiaoming",
			ToList: []string{"ID_xiaowang"},
		},
	}

	chatRecordService.EXPECT().Read(gomock.Eq(uint64(267)), gomock.Eq(uint64(10))).Return(wecomRecords, nil).Times(1)
	recordTransformer.EXPECT().Transform(gomock.Any()).Return(nil, io.ErrUnexpectedEOF).Times(1)

	// When
	_, _, err := readerAdapter.Read(pagination.NewPageToken(267), 10)
	if err == nil {
		t.Errorf("error should happen here, expected: %v, actual: %v", io.ErrUnexpectedEOF, err)
		return
	}
}
