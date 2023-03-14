package use_case

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
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
