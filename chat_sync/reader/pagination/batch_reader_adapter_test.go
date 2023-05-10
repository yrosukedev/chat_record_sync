package pagination

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"io"
	"reflect"
	"testing"
)

func TestFetchPageToken_nil(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReader := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(nil, nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Nil(), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, NewPageToken(0), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(0))).Return(nil).Times(1)

	// When
	batchReader.Read()
}

func TestFetchPageToken_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(0), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(0)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, NewPageToken(0), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(0))).Return(nil).Times(1)

	// When
	batchReaderAdapter.Read()
}

func TestFetchPageToken_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(1), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(1)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, NewPageToken(1), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(1))).Return(nil).Times(1)

	// When
	batchReaderAdapter.Read()
}

func TestFetchPageToken_many(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(3), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(3)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, NewPageToken(3), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(3))).Return(nil).Times(1)

	// When
	batchReaderAdapter.Read()
}

func TestFetchPageToken_ALot(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(10500), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10500)), gomock.Eq(pageSize)).Return([]*business.ChatRecord{}, NewPageToken(10500), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(10500))).Return(nil).Times(1)

	// When
	batchReaderAdapter.Read()
}

func TestFetchPageToken_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(0), io.ErrUnexpectedEOF).Times(1)
	paginatedReader.EXPECT().Read(gomock.Any(), gomock.Any()).Times(0)
	paginationStorage.EXPECT().Set(gomock.Any()).Times(0)

	// When
	_, err := batchReaderAdapter.Read()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrUnexpectedEOF, err)
	}
}

func TestForwardResults_zeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(123456), nil).Times(1)

	var records []*business.ChatRecord
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(123456)), gomock.Eq(pageSize)).Return(records, NewPageToken(897654), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(897654))).Return(nil).Times(1)

	// When
	forwardingResults, err := batchReaderAdapter.Read()
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
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(123456), nil).Times(1)

	records := []*business.ChatRecord{
		{},
	}
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(123456)), gomock.Eq(pageSize)).Return(records, NewPageToken(2234567), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(2234567))).Return(nil).Times(1)

	// When
	forwardingResults, err := batchReaderAdapter.Read()
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
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)

	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(records, NewPageToken(678934), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(678934))).Return(nil).Times(1)

	// When
	forwardingResults, err := batchReaderAdapter.Read()
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
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)

	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(nil, NewPageToken(0), io.ErrShortBuffer).Times(1)

	paginationStorage.EXPECT().Set(gomock.Any()).Times(0)

	// When
	_, err := batchReaderAdapter.Read()
	if err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
}

func TestForwardResults_EOF(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)

	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(nil, NewPageToken(3467), io.EOF).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(3467))).Return(nil).Times(0)

	// When
	_, err := batchReaderAdapter.Read()
	if err != io.EOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.EOF, err)
	}
}

func TestForwardResults_ManyEOF(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	expectedRecords := []*business.ChatRecord{
		{},
		{},
		{},
	}

	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(nil, NewPageToken(20), io.EOF).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(20))).Return(nil).Times(0)

	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(expectedRecords, NewPageToken(30), nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(30))).Return(nil).Times(1)

	paginationStorage.EXPECT().Get().Return(NewPageToken(30), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(30)), gomock.Eq(pageSize)).Return(nil, NewPageToken(40), io.EOF).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(40))).Return(nil).Times(0)

	// When & Then
	_, err := batchReaderAdapter.Read()
	assert.ErrorIs(t, err, io.EOF)

	records, err := batchReaderAdapter.Read()
	if assert.NoError(t, err) {
		assert.Equal(t, expectedRecords, records)
	}

	_, err = batchReaderAdapter.Read()
	assert.ErrorIs(t, err, io.EOF)

	_, err = batchReaderAdapter.Read()
	assert.ErrorIs(t, err, io.EOF)
}

func TestUpdatePageToken_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	// Then
	paginationStorage.EXPECT().Get().Return(NewPageToken(10), nil).Times(1)

	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(10)), gomock.Eq(pageSize)).Return(records, NewPageToken(678934), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(678934))).Return(io.ErrUnexpectedEOF).Times(1)

	// When
	_, err := batchReaderAdapter.Read()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrUnexpectedEOF, err)
	}
}

func TestDetermineEnd_requestPageSizeEqualToResponsePageSize(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(5)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	records := []*business.ChatRecord{
		{},
		{},
		{},
		{},
		{},
	}

	// Then
	givenPaginationStoragePageTokens(paginationStorage, []*PageToken{NewPageToken(100), NewPageToken(231456)})

	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(100)), gomock.Eq(pageSize)).Return(records, NewPageToken(231456), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(231456)), gomock.Eq(pageSize)).Return(records, NewPageToken(901234), nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(231456))).Return(nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(901234))).Return(nil).Times(1)

	// When

	// 1st reading operation
	forwardingResults, err := batchReaderAdapter.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}

	// 2nd reading operation
	forwardingResults, err = batchReaderAdapter.Read()
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
	paginatedReader := NewMockPaginatedReader(ctrl)
	paginationStorage := NewMockPaginationStorage(ctrl)
	pageSize := uint64(10)
	batchReaderAdapter := NewBatchReaderAdapter(paginatedReader, paginationStorage, pageSize)

	records := []*business.ChatRecord{
		{},
		{},
		{},
		{},
		{},
		{},
	}

	// Then
	givenPaginationStoragePageTokens(paginationStorage, []*PageToken{NewPageToken(2000), NewPageToken(547612)})

	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(2000)), gomock.Eq(pageSize)).Return(records, NewPageToken(547612), nil).Times(1)
	paginatedReader.EXPECT().Read(gomock.Eq(NewPageToken(547612)), gomock.Eq(pageSize)).Return(records, NewPageToken(657831), nil).Times(0)

	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(547612))).Return(nil).Times(1)
	paginationStorage.EXPECT().Set(gomock.Eq(NewPageToken(657831))).Return(nil).Times(0)

	// When

	// 1st reading operation
	forwardingResults, err := batchReaderAdapter.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}

	// 2nd reading operation
	_, err = batchReaderAdapter.Read()
	if err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %v, actual: %v", io.EOF, err)
	}
}

func givenPaginationStoragePageTokens(paginationStorage *MockPaginationStorage, pageTokens []*PageToken) {
	idx := 0
	paginationStorage.
		EXPECT().
		Get().
		DoAndReturn(func() (*PageToken, error) {
			if idx < len(pageTokens) {
				defer func() { idx += 1 }()
				return pageTokens[idx], nil
			}
			return NewPageToken(0), fmt.Errorf("page token out of range, index: %v, length of token array: %v", idx, len(pageTokens))
		}).
		AnyTimes()
}
