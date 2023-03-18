package buffer_reader

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
	"testing"
)

func TestBufferSize_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	records := []*business.ChatRecord{
		{},
	}
	bufferedReader.EXPECT().Read().Return(records, nil).Times(10)

	// When
	for i := 0; i < 10; i++ {
		_, _ = readerAdapter.Read()
	}
}

func TestBufferSize_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	var records []*business.ChatRecord
	bufferedReader.EXPECT().Read().Return(records, nil).Times(10)

	// When
	for i := 0; i < 10; i++ {
		if _, err := readerAdapter.Read(); err != io.EOF {
			t.Errorf("end should happen here, expected: %+v, actual: %+v", io.EOF, err)
		}
	}
}

func TestBufferSize_greaterThanOne(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	records := []*business.ChatRecord{
		{
			MsgId: "1",
		},
		{
			MsgId: "2",
		},
		{
			MsgId: "3",
		},
	}
	bufferedReader.EXPECT().Read().Return(records, nil).Times(1)

	// When
	expectReaderToReadRecords(t, readerAdapter, records)
}

func TestBufferSize_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	bufferedReader.EXPECT().Read().Return(nil, io.ErrShortBuffer).Times(1)

	// When
	if _, err := readerAdapter.Read(); err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
}

func TestRefill_zeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	recordsGroup1 := []*business.ChatRecord{
		{
			MsgId: "1",
		},
		{
			MsgId: "2",
		},
		{
			MsgId: "3",
		},
	}
	var recordsGroup2 []*business.ChatRecord
	givenRecordsToRefill(bufferedReader, [][]*business.ChatRecord{recordsGroup1, recordsGroup2})

	// When
	expectReaderToReadRecords(t, readerAdapter, append(recordsGroup1))
	if _, err := readerAdapter.Read(); err != io.EOF {
		t.Errorf("end should happen here, expected: %+v, actual: %+v", io.EOF, err)
	}
}

func givenRecordsToRefill(bufferedReader *MockChatRecordBufferedReader, recordsGroups [][]*business.ChatRecord) {
	groupIdx := 0
	bufferedReader.EXPECT().Read().DoAndReturn(func() ([]*business.ChatRecord, error) {
		defer func() { groupIdx += 1 }()
		return recordsGroups[groupIdx], nil
	}).Times(len(recordsGroups))
}

func expectReaderToReadRecords(t *testing.T, reader *ChatRecordBufferedReaderAdapter, records []*business.ChatRecord) {
	for _, expected := range records {
		actual, err := reader.Read()
		if err != nil {
			t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
		}

		if expected != actual {
			t.Errorf("records are not matched, expected: %+v, actual: %+v", expected, actual)
		}
	}
}
