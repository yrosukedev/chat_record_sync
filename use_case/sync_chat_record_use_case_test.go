package use_case

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"golang.org/x/exp/slices"
	"io"
	"reflect"
	"testing"
)

func TestZeroRecord(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	reader.EXPECT().Read().Times(1).Return(nil, io.EOF)

	// Then
	writer.EXPECT().Write(gomock.Any()).Times(0)

	// When
	useCase.Run(ctx)
}

func TestOneRecord(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
	}
	givenRecordsToRead(reader, records)

	// Then
	writer.
		EXPECT().
		Write(gomock.Eq(records[0])).
		Times(1)

	// When
	useCase.Run(ctx)
}

func TestManyRecords(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	givenRecordsToRead(reader, records)

	// Then
	expectRecordsToWrite(writer, records)

	// When
	useCase.Run(ctx)
}

func TestReaderError_beforeReading(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	encounterErrorWhileReadingRecords(reader, records, []int{0}, io.ErrClosedPipe)

	// Then
	writer.EXPECT().Write(records[1]).Times(1)

	// When
	useCase.Run(ctx)
}

func TestReaderError_whileReading(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	encounterErrorWhileReadingRecords(reader, records, []int{1}, io.ErrClosedPipe)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[0])).Times(1)
	writer.EXPECT().Write(gomock.Eq(records[2])).Times(1)

	// When
	useCase.Run(ctx)
}

func TestWriterError_firstRecord(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	givenRecordsToRead(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[0])).Times(1).Return(io.ErrShortWrite)
	writer.EXPECT().Write(gomock.Eq(records[1])).Times(1).Return(nil)
	writer.EXPECT().Write(gomock.Eq(records[2])).Times(1).Return(nil)

	// When
	useCase.Run(ctx)
}

func TestWriterError_recordInTheMiddle(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	givenRecordsToRead(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[0])).Times(1).Return(nil)
	writer.EXPECT().Write(gomock.Eq(records[1])).Times(1).Return(io.ErrNoProgress)
	writer.EXPECT().Write(gomock.Eq(records[2])).Times(1).Return(nil)

	// When
	useCase.Run(ctx)
}

func TestAccumulateErrors_readerError(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	encounterErrorWhileReadingRecords(reader, records, []int{0, 2}, io.ErrClosedPipe)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[1])).Times(1).Return(nil)

	// When
	errs := useCase.Run(ctx)
	expectedErrs := []*SyncError{
		NewReaderError(io.ErrClosedPipe),
		NewReaderError(io.ErrClosedPipe),
	}
	if !reflect.DeepEqual(errs, expectedErrs) {
		t.Errorf("accumulated errors are not equal, expected: %+v, actual: %+v", expectedErrs, errs)
	}
}

func TestAccumulatedErrors_writeError(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*business.ChatRecord{
		&business.ChatRecord{},
		&business.ChatRecord{},
		&business.ChatRecord{},
	}
	givenRecordsToRead(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[0])).Times(1).Return(nil)
	writer.EXPECT().Write(gomock.Eq(records[1])).Times(1).Return(io.ErrShortBuffer)
	writer.EXPECT().Write(gomock.Eq(records[2])).Times(1).Return(io.ErrShortWrite)

	// When
	errs := useCase.Run(ctx)
	expectedErrs := []*SyncError{
		NewWriterError(io.ErrShortBuffer, records[1]),
		NewWriterError(io.ErrShortWrite, records[2]),
	}
	if !reflect.DeepEqual(errs, expectedErrs) {
		t.Errorf("accumulated errors are not equal, expected: %+v, actual: %+v", expectedErrs, errs)
	}
}

func expectRecordsToWrite(writer *MockChatRecordWriter, records []*business.ChatRecord) {
	for _, r := range records {
		writer.
			EXPECT().
			Write(gomock.Eq(r)).
			Times(1)
	}
}

func givenRecordsToRead(reader *MockChatRecordReader, records []*business.ChatRecord) {
	idx := 0
	reader.
		EXPECT().
		Read().
		DoAndReturn(func() (*business.ChatRecord, error) {
			if idx >= len(records) {
				return nil, io.EOF
			}
			defer func() { idx += 1 }()
			return records[idx], nil
		}).
		AnyTimes()
}

func encounterErrorWhileReadingRecords(reader *MockChatRecordReader, records []*business.ChatRecord, errIdxs []int, err error) {
	idx := 0
	reader.
		EXPECT().
		Read().
		DoAndReturn(func() (*business.ChatRecord, error) {
			if idx >= len(records) {
				return nil, io.EOF
			}
			defer func() { idx += 1 }()
			if slices.Contains(errIdxs, idx) {
				return nil, err
			}
			return records[idx], nil
		}).
		AnyTimes()
}
