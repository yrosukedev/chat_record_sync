package buffer

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"io"
	"testing"
)

func TestBufferSize_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

	// Then
	records := []*business.ChatRecord{
		{},
	}
	batchReader.EXPECT().Read().Return(records, nil).Times(10)

	// When
	for i := 0; i < 10; i++ {
		_, _ = readerAdapter.Read()
	}
}

func TestBufferSize_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

	// Then
	var records []*business.ChatRecord
	batchReader.EXPECT().Read().Return(records, nil).Times(10)

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
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

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
	batchReader.EXPECT().Read().Return(records, nil).Times(1)

	// When
	expectReaderToReadRecords(t, readerAdapter, records)
}

func TestBufferSize_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

	// Then
	batchReader.EXPECT().Read().Return(nil, io.ErrShortBuffer).Times(1)

	// When
	if _, err := readerAdapter.Read(); err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
}

func TestRefill_zeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

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
	givenRecordsToRefill(batchReader, [][]*business.ChatRecord{recordsGroup1, recordsGroup2})

	// When
	expectReaderToReadRecords(t, readerAdapter, recordsGroup1)
	if _, err := readerAdapter.Read(); err != io.EOF {
		t.Errorf("end should happen here, expected: %+v, actual: %+v", io.EOF, err)
	}
}

func TestRefill_oneTime(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

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
	recordsGroup2 := []*business.ChatRecord{
		{
			MsgId: "4",
		},
	}
	givenRecordsToRefill(batchReader, [][]*business.ChatRecord{recordsGroup1, recordsGroup2})

	// When
	expectReaderToReadRecords(t, readerAdapter, append(recordsGroup1, recordsGroup2...))
}

func TestRefill_manyTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

	// Then
	recordsGroup1 := []*business.ChatRecord{
		{
			MsgId: "1",
		},
		{
			MsgId: "2",
		},
	}
	recordsGroup2 := []*business.ChatRecord{
		{
			MsgId: "3",
		},
		{
			MsgId: "4",
		},
		{
			MsgId: "5",
		},
	}
	recordsGroup3 := []*business.ChatRecord{
		{
			MsgId: "6",
		},
	}
	recordsGroup4 := []*business.ChatRecord{
		{
			MsgId: "7",
		},
		{
			MsgId: "8",
		},
	}
	givenRecordsToRefill(
		batchReader,
		[][]*business.ChatRecord{
			recordsGroup1,
			recordsGroup2,
			recordsGroup3,
			recordsGroup4,
		},
	)

	// When
	expectReaderToReadRecords(t, readerAdapter, append(append(append(recordsGroup1, recordsGroup2...), recordsGroup3...), recordsGroup4...))
}

func TestRefill_error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

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

	groupIdx := 0
	batchReader.EXPECT().Read().DoAndReturn(func() ([]*business.ChatRecord, error) {
		defer func() { groupIdx += 1 }()
		if groupIdx == 0 {
			return records, nil
		}
		return nil, io.ErrUnexpectedEOF
	}).Times(2)

	// When
	expectReaderToReadRecords(t, readerAdapter, records)
	if _, err := readerAdapter.Read(); err != io.ErrUnexpectedEOF {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrUnexpectedEOF, err)
	}
}

func TestReader_Read_ManyEOF(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	batchReader := NewMockBatchReader(ctrl)
	readerAdapter := NewReader(batchReader)

	recordsGroup1 := []*business.ChatRecord{
		{
			MsgId: "1",
		},
		{
			MsgId: "2",
		},
	}
	recordsGroup2 := []*business.ChatRecord{
		{
			MsgId: "3",
		},
		{
			MsgId: "4",
		},
		{
			MsgId: "5",
		},
	}
	recordsGroup3 := []*business.ChatRecord{
		{
			MsgId: "6",
		},
	}
	recordsGroup4 := []*business.ChatRecord{
		{
			MsgId: "7",
		},
		{
			MsgId: "8",
		},
	}
	givenRecordsOrErrorsToRefill(batchReader, []recordsOrErr{
		&recordsWrapper{
			_records: recordsGroup1,
		},
		&errWrapper{
			_err: io.EOF,
		},
		&recordsWrapper{
			_records: recordsGroup2,
		},
		&errWrapper{
			_err: io.EOF,
		},
		&recordsWrapper{
			_records: recordsGroup3,
		},
		&errWrapper{
			_err: io.EOF,
		},
		&recordsWrapper{
			_records: recordsGroup4,
		},
		&errWrapper{
			_err: io.EOF,
		},
	})

	// When
	expectReaderToReadRecordOrError(t, readerAdapter, []use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(recordsGroup1[0]),
		use_case.NewRecordOrErrorWithRecord(recordsGroup1[1]),
		use_case.NewRecordOrErrorWithError(io.EOF),
		use_case.NewRecordOrErrorWithRecord(recordsGroup2[0]),
		use_case.NewRecordOrErrorWithRecord(recordsGroup2[1]),
		use_case.NewRecordOrErrorWithRecord(recordsGroup2[2]),
		use_case.NewRecordOrErrorWithError(io.EOF),
		use_case.NewRecordOrErrorWithRecord(recordsGroup3[0]),
		use_case.NewRecordOrErrorWithError(io.EOF),
		use_case.NewRecordOrErrorWithRecord(recordsGroup4[0]),
		use_case.NewRecordOrErrorWithRecord(recordsGroup4[1]),
		use_case.NewRecordOrErrorWithError(io.EOF),
	})
}
