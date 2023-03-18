package paginated_reader

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
	"reflect"
	"testing"
)

func TestFetchPageToken_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(0), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(0)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, PageToken(0), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(0))).Return(nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(1), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(1)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, PageToken(1), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(1))).Return(nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_many(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(3), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(3)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, PageToken(3), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(3))).Return(nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_ALot(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(10500), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(10500)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, PageToken(10500), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(10500))).Return(nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(0), io.ErrUnexpectedEOF).Times(1)
	bufferedReader.EXPECT().Read(gomock.Any(), gomock.Any()).Times(0)
	paginationStorage.EXPECT().Set(gomock.Any()).Times(0)

	// When
	_, err := paginatedReader.Read()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrUnexpectedEOF, err)
	}
}

func TestForwardResults_zeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(123456), nil).Times(1)

	var records []*business.ChatRecord
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(123456)), gomock.Eq(pageSize)).Return(records, PageToken(897654), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(897654))).Return(nil).Times(1)

	// When
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}
}

func TestForwardResults_oneRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(123456), nil).Times(1)

	records := []*business.ChatRecord{
		{},
	}
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(123456)), gomock.Eq(pageSize)).Return(records, PageToken(2234567), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(2234567))).Return(nil).Times(1)

	// When
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}
}

func TestForwardResults_manyRecords(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(10), nil).Times(1)

	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(10)), gomock.Eq(pageSize)).Return(records, PageToken(678934), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(678934))).Return(nil).Times(1)

	// When
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}
}

func TestForwardResults_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(10), nil).Times(1)

	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(10)), gomock.Eq(pageSize)).Return(nil, PageToken(0), io.ErrShortBuffer).Times(1)

	paginationStorage.EXPECT().Set(gomock.Any()).Times(0)

	// When
	_, err := paginatedReader.Read()
	if err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
}

func TestUpdatePageToken_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(10), nil).Times(1)

	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(10)), gomock.Eq(pageSize)).Return(records, PageToken(678934), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(678934))).Return(io.ErrUnexpectedEOF).Times(1)

	// When
	_, err := paginatedReader.Read()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrUnexpectedEOF, err)
	}
}

func TestDetermineEnd_requestPageSizeEqualToResponsePageSize(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(5)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	records := []*business.ChatRecord{
		{},
		{},
		{},
		{},
		{},
	}

	// Then
	givenPaginationStoragePageTokens(paginationStorage, []PageToken{PageToken(100), PageToken(231456)})

	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(100)), gomock.Eq(pageSize)).Return(records, PageToken(231456), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(231456)), gomock.Eq(pageSize)).Return(records, PageToken(901234), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(231456))).Return(nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(901234))).Return(nil).Times(1)

	// When

	// 1st reading operation
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}

	// 2nd reading operation
	forwardingResults, err = paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}
}

func TestDetermineEnd_requestPageSizeGreaterThanResponsePageSize(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	pageSize := uint64(10)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage, pageSize)

	records := []*business.ChatRecord{
		{},
		{},
		{},
		{},
		{},
		{},
	}

	// Then
	givenPaginationStoragePageTokens(paginationStorage, []PageToken{PageToken(2000), PageToken(547612)})

	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(2000)), gomock.Eq(pageSize)).Return(records, PageToken(547612), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(547612)), gomock.Eq(pageSize)).Return(records, PageToken(657831), nil).Times(0)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(547612))).Return(nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(657831))).Return(nil).Times(0)

	// When

	// 1st reading operation
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}

	// 2nd reading operation
	_, err = paginatedReader.Read()
	if err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %v, actual: %v", io.EOF, err)
	}
}

func givenPaginationStoragePageTokens(paginationStorage *MockChatRecordPaginationStorage, pageTokens []PageToken) {
	idx := 0
	paginationStorage.
		EXPECT().
		Get().
		DoAndReturn(func() (PageToken, error) {
			if idx < len(pageTokens) {
				defer func() { idx += 1 }()
				return pageTokens[idx], nil
			}
			return PageToken(0), fmt.Errorf("page token out of range, index: %v, length of token array: %v", idx, len(pageTokens))
		}).
		AnyTimes()
}
