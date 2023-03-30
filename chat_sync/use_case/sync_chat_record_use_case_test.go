package use_case

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
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
	GivenRecordsToRead(reader, records)

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
	GivenRecordsToRead(reader, records)

	// Then
	ExpectRecordsToWrite(writer, records)

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

	records := []*RecordOrError{
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	writer.EXPECT().Write(records[1].Record).Times(1)

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

	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[0].Record)).Times(1)
	writer.EXPECT().Write(gomock.Eq(records[2].Record)).Times(1)

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
	GivenRecordsToRead(reader, records)

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
	GivenRecordsToRead(reader, records)

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

	records := []*RecordOrError{
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrNoProgress),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[1].Record)).Times(1).Return(nil)

	// When
	errs := useCase.Run(ctx)
	expectedErrs := []*SyncError{
		NewReaderError(io.ErrClosedPipe),
		NewReaderError(io.ErrNoProgress),
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
	GivenRecordsToRead(reader, records)

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

func TestAccumulateErrors_readAndWriteErrors(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	records := []*RecordOrError{
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	writer.EXPECT().Write(gomock.Eq(records[2].Record)).Times(1).Return(io.ErrNoProgress)

	// When
	errs := useCase.Run(ctx)
	expectedErrs := []*SyncError{
		NewReaderError(io.ErrUnexpectedEOF),
		NewReaderError(io.ErrShortBuffer),
		NewWriterError(io.ErrNoProgress, records[2].Record),
	}
	if !reflect.DeepEqual(errs, expectedErrs) {
		t.Errorf("accumulated errors are not equal, expected: %+v, actual: %+v", expectedErrs, errs)
	}
}
